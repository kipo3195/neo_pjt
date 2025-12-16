package handler

import (
	"io"
	"log"
	"org/internal/application/orchestrator"
	"org/internal/consts"
	"org/internal/delivery/adapter"
	"org/internal/delivery/dto/org"
	commonConsts "org/pkg/consts"
	"org/pkg/response"

	"github.com/gin-gonic/gin"
)

type OrgBatchServiceHandler struct {
	svc *orchestrator.OrgBatchService
}

func NewOrgBatchServiceHandler(svc *orchestrator.OrgBatchService) *OrgBatchServiceHandler {
	return &OrgBatchServiceHandler{
		svc: svc,
	}
}

func (r *OrgBatchServiceHandler) RegistOrgBatchData(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	orgCode := c.Request.FormValue("org_code")
	log.Println("[RegistOrgBatchData] orgCode :", orgCode)

	// 요청 org code가 없는 경우
	if orgCode == "" {
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		return
	}

	// 파일 데이터 추출
	orgFile, err := c.FormFile(consts.ORG_FILE)
	if err != nil {
		// 파일
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_103, commonConsts.E_103_MSG)
		return
	}

	// 파일 처리 []byte -> dto
	file, err := orgFile.Open()
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

	req := org.RegistOrgBatchRequest{
		OrgFile:     &fileBytes,
		OrgFileName: orgFile.Filename,
		OrgCode:     orgCode,
	}

	input := adapter.MakeRegistOrgBatchInput(req.OrgFile, req.OrgFileName, req.OrgCode)
	err = r.svc.RegistOrgBatch(ctx, input)

	if err != nil {
		log.Println("error : ", err)
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		response.SendSuccess(c, "success")
	}

}
