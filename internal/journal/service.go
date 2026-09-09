package journal

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"

	djv1 "dj/gen/go/dj/v1"
	"dj/gen/go/dj/v1/djv1connect"

	"dj/internal/config"
	"dj/internal/logs"
)

// Service implements journal use cases.
type Service struct {
	store  JournalStore
	config config.Provider
	log    logs.Logger
	now    func() time.Time
}

// NewService returns a journal application service.
func NewService(store JournalStore, cfg config.Provider, log logs.Logger) *Service {
	return &Service{
		store:  store,
		config: cfg,
		log:    log,
		now:    time.Now,
	}
}

func (s *Service) StartDay(ctx context.Context) (*djv1.StartDayResponse, error) {
	path, created, err := s.store.CreateToday(ctx)
	if err != nil {
		return nil, err
	}
	return &djv1.StartDayResponse{
		FilePath: path,
		Created:  created,
	}, nil
}

func (s *Service) AddEntry(ctx context.Context, text string) (*djv1.AddEntryResponse, error) {
	path, err := s.store.TodaysPath(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.store.AppendEntry(ctx, path, text); err != nil {
		return nil, err
	}

	return &djv1.AddEntryResponse{FilePath: path}, nil
}

func (s *Service) AddBullet(ctx context.Context, text string) (*djv1.AddBulletResponse, error) {
	path, err := s.store.TodaysPath(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.store.AppendBullet(ctx, path, text); err != nil {
		return nil, err
	}

	return &djv1.AddBulletResponse{FilePath: path}, nil
}

func (s *Service) GetDay(ctx context.Context, req *djv1.GetDayRequest) (*djv1.GetDayResponse, error) {
	date, err := s.resolveDate(req)
	if err != nil {
		return nil, err
	}

	content, path, err := s.store.ReadDay(ctx, date)
	if err != nil {
		return nil, err
	}

	return &djv1.GetDayResponse{
		Day: &djv1.JournalDay{
			Date:     timestamppb.New(date),
			FilePath: path,
			Content:  content,
		},
	}, nil
}

func (s *Service) MigrateFolders(ctx context.Context, root string, dryRun bool) (*djv1.MigrateFoldersResponse, error) {
	renames, err := s.store.MigrateFolders(ctx, root, dryRun)
	if err != nil {
		return nil, err
	}
	return &djv1.MigrateFoldersResponse{Renames: renames}, nil
}

func (s *Service) ListMonthFolders(ctx context.Context) (*djv1.ListMonthFoldersResponse, error) {
	folders, err := s.store.ListMonthFolders(ctx)
	if err != nil {
		return nil, err
	}
	return &djv1.ListMonthFoldersResponse{FolderNames: folders}, nil
}

func (s *Service) ListFiles(ctx context.Context, monthFolder string) (*djv1.ListFilesResponse, error) {
	files, err := s.store.ListFiles(ctx, monthFolder)
	if err != nil {
		return nil, err
	}
	return &djv1.ListFilesResponse{FileNames: files}, nil
}

func (s *Service) ReadFile(ctx context.Context, monthFolder, fileName string) (*djv1.ReadFileResponse, error) {
	content, err := s.store.ReadFile(ctx, monthFolder, fileName)
	if err != nil {
		return nil, err
	}
	return &djv1.ReadFileResponse{Content: content}, nil
}

func (s *Service) JournalDirectory() string {
	return s.config.JournalDirectory()
}

func (s *Service) resolveDate(req *djv1.GetDayRequest) (time.Time, error) {
	switch selector := req.Selector.(type) {
	case *djv1.GetDayRequest_Today:
		return s.now(), nil
	case *djv1.GetDayRequest_Yesterday:
		return s.now().AddDate(0, 0, -1), nil
	case *djv1.GetDayRequest_Date:
		return selector.Date.AsTime(), nil
	case *djv1.GetDayRequest_DateString:
		date, err := time.Parse("01/02/2006", selector.DateString)
		if err != nil {
			return time.Time{}, fmt.Errorf("%w: %s", ErrInvalidDate, err)
		}
		return date, nil
	default:
		return time.Time{}, ErrInvalidDate
	}
}

// ConnectHandler adapts Service to the Connect JournalService handler interface.
type ConnectHandler struct {
	svc *Service
}

// NewConnectHandler returns a Connect RPC handler for JournalService.
func NewConnectHandler(svc *Service) *ConnectHandler {
	return &ConnectHandler{svc: svc}
}

var _ djv1connect.JournalServiceHandler = (*ConnectHandler)(nil)

func (h *ConnectHandler) StartDay(ctx context.Context, req *connect.Request[djv1.StartDayRequest]) (*connect.Response[djv1.StartDayResponse], error) {
	resp, err := h.svc.StartDay(ctx)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (h *ConnectHandler) AddEntry(ctx context.Context, req *connect.Request[djv1.AddEntryRequest]) (*connect.Response[djv1.AddEntryResponse], error) {
	resp, err := h.svc.AddEntry(ctx, req.Msg.GetText())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (h *ConnectHandler) AddBullet(ctx context.Context, req *connect.Request[djv1.AddBulletRequest]) (*connect.Response[djv1.AddBulletResponse], error) {
	resp, err := h.svc.AddBullet(ctx, req.Msg.GetText())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (h *ConnectHandler) GetDay(ctx context.Context, req *connect.Request[djv1.GetDayRequest]) (*connect.Response[djv1.GetDayResponse], error) {
	resp, err := h.svc.GetDay(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (h *ConnectHandler) MigrateFolders(ctx context.Context, req *connect.Request[djv1.MigrateFoldersRequest]) (*connect.Response[djv1.MigrateFoldersResponse], error) {
	resp, err := h.svc.MigrateFolders(ctx, req.Msg.GetRoot(), req.Msg.GetDryRun())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (h *ConnectHandler) ListMonthFolders(ctx context.Context, req *connect.Request[djv1.ListMonthFoldersRequest]) (*connect.Response[djv1.ListMonthFoldersResponse], error) {
	resp, err := h.svc.ListMonthFolders(ctx)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (h *ConnectHandler) ListFiles(ctx context.Context, req *connect.Request[djv1.ListFilesRequest]) (*connect.Response[djv1.ListFilesResponse], error) {
	resp, err := h.svc.ListFiles(ctx, req.Msg.GetMonthFolder())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (h *ConnectHandler) ReadFile(ctx context.Context, req *connect.Request[djv1.ReadFileRequest]) (*connect.Response[djv1.ReadFileResponse], error) {
	resp, err := h.svc.ReadFile(ctx, req.Msg.GetMonthFolder(), req.Msg.GetFileName())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
