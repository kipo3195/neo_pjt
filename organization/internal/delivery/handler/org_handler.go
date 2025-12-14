package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"org/internal/application/usecase"
	"org/internal/delivery/adapter"
	"org/internal/delivery/dto/org"
	commonConsts "org/pkg/consts"
	"org/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type OrgHandler struct {
	usecase usecase.OrgUsecase
}

func NewOrgHandler(usecase usecase.OrgUsecase) *OrgHandler {
	return &OrgHandler{
		usecase: usecase,
	}
}

// 조직도 전체 조회
func (h *OrgHandler) GetOrgHash(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var req = org.GetOrgHashRequest{
		// 배열의 형태로 받음. org가 하나 이상일 수도 있기 때문.
		OrgHash: c.QueryArray("orgHash"),
	}

	if len(req.OrgHash) == 0 {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}

	// usecase 호출
	data, err := h.usecase.GetOrgHash(ctx, req)

	// response.
	if err == nil {
		response.SendSuccess(c, data)
	} else {
		// http status code 400
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	}
}

func (h *OrgHandler) GetOrgData(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var req = org.GetOrgDataRequest{
		OrgCode: c.Query("orgCode"),
		Type:    c.Query("type"),
		OrgHash: c.Query("orgHash"),
	}

	if len(req.OrgCode) == 0 {
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_108, commonConsts.E_108_MSG)
		return
	}
	// usecase 호출
	fileName, data, err := h.usecase.GetOrgData(ctx, req)

	// response.
	if fileName != "" {
		orgCode := req.OrgCode
		// http status code 200
		// w.Header().Set("Content-Type", "application/octet-stream")
		// w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.zip"`, orgCode+"_"+file)) // 요청한 org code + 최신 hash
		// w.Write(data.([]byte))
		// // 전송 헤더의 순서가 영향을 미침 - 파일명 적용이 안됨.
		// w.WriteHeader(http.StatusOK)

		// "." 기준으로 나누기
		parts := strings.SplitN(fileName, ".", 2) // 최대 2개만 분리
		// 앞부분만 사용
		date := parts[0]

		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, orgCode+"_"+date+".zip"))
		c.Data(http.StatusOK, "application/octet-stream", data.([]byte))

		// interface{} → *os.File 로 변환
		// if realFile, ok := data.(*os.File); ok {
		// 	fmt.Println("변환 성공:", realFile.Name())
		// 	response.SendFileStream(c, realFile, orgCode+"_"+fileName+".zip", "")
		// } else {
		// 	fmt.Println("변환 실패")
		// 	response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
		// }
		return

	} else if err != nil {
		// http status code 400
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		response.SendSuccess(c, data)
	}
}

func (h *OrgHandler) RegistOrgBatch(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	// 파일 데이터 추출
	orgFile, err := c.FormFile("org_file")
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
	}

	input := adapter.MakeRegistOrgBatchInput(req.OrgFile, req.OrgFileName)
	err = h.usecase.RegistOrgBatch(ctx, input)

	if err != nil {
		log.Println("error : ", err)
		response.SendError(c, commonConsts.SERVER_ERROR, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	} else {
		response.SendSuccess(c, "success")
	}

}

func (h *OrgHandler) CreateOrgFile(c *gin.Context) {

	// context 생성
	ctx := c.Request.Context()

	// request 데이터 파싱 header, body -> dto
	var req = org.CreateOrgFileRequest{}

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

	// usecase 호출
	data, err := h.usecase.CreateOrgFile(ctx, req)

	if err == nil {
		// http status code 200
		response.SendSuccess(c, data)
	} else {
		// http status code 500
		response.SendError(c, commonConsts.BAD_REQUEST, commonConsts.ERROR, commonConsts.E_500, commonConsts.E_500_MSG)
	}

	// response.

}
