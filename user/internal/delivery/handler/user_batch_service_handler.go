package handler

import (
	"io"
	"log"
	"user/internal/application/orchestrator"
	"user/internal/consts"
	"user/internal/delivery/adapter"
	"user/internal/delivery/dto/userBatch"
	commonConsts "user/pkg/consts"
	"user/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserBatchServiceHandler struct {
	svc *orchestrator.UserBatchService
}

func NewUserBatchServiceHandler(svc *orchestrator.UserBatchService) *UserBatchServiceHandler {
	return &UserBatchServiceHandler{
		svc: svc,
	}
}

func (r *UserBatchServiceHandler) RegistUserDetailData(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	orgCode := c.Request.FormValue("org_code")
	log.Println("[RegistUserDetailData] orgCode :", orgCode)

	// 요청 org code가 없는 경우
	if orgCode == "" {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	// 파일 데이터 추출
	userDetailFile, err := c.FormFile(consts.USER_DETAIL_FILE)
	if err != nil {
		// 파일
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 파일 처리 []byte -> dto
	file, err := userDetailFile.Open()
	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	// Open() 이 에러를 반환하기 전에 defer Close()를 걸면 안 된다
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	req := userBatch.RegistUserBatchRequest{
		File:     &fileBytes,
		FileName: userDetailFile.Filename,
		OrgCode:  orgCode,
	}

	input := adapter.MakeRegistUserDetailBatchInput(req.File, req.FileName, req.OrgCode)
	err = r.svc.RegisterUserDetailBatch(ctx, input)

	if err != nil {
		log.Println("error : ", err)
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		response.SendSuccess(c, "success")
	}

}
