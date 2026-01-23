package handler

import (
	"encoding/json"
	"file/internal/application/usecase"
	"file/internal/delivery/adapter"
	"file/internal/delivery/dto/fileUrl"
	"file/internal/delivery/util"
	commonConsts "file/pkg/consts"
	"file/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type FileUrlHandler struct {
	usecase usecase.FileUrlUsecase
}

func NewFileUrlHandler(usecase usecase.FileUrlUsecase) *FileUrlHandler {
	return &FileUrlHandler{
		usecase: usecase,
	}
}

func (r *FileUrlHandler) CreateFileUrl(c *gin.Context) {

	ctx := c.Request.Context()

	// 사용자 정보 파싱
	reqUserHash := util.GetUserHashByAccessToken(c)
	if reqUserHash == "" {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_110, commonConsts.E_110_MSG)
		return
	}

	var req fileUrl.CreateFileUrlRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 필수 데이터 검증
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	input := adapter.MakeCreateFileUrlInput(reqUserHash, req.EventType, req.Org, req.FileInfo)

	output, err := r.usecase.CreateFileUrl(ctx, input)

	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	fileUrlInfo := make([]fileUrl.FileUrlInfoDto, 0)

	for _, f := range output.FileUrlInfo {

		temp := fileUrl.FileUrlInfoDto{
			FileName: f.FileName,
			Url:      f.Url,
		}

		fileUrlInfo = append(fileUrlInfo, temp)

	}

	res := fileUrl.CreateFileUrlResponse{
		TransactionId: output.TransactionId,
		FileUrlInfo:   fileUrlInfo,
	}

	response.SendSuccess(c, res)

}
