package usecases

import (
	"archive/zip"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"org/consts"
	clDto "org/dto/client"
	svDto "org/dto/server"
	"org/entities"
	"org/models"
	"org/repositories"
	"os"
	"strings"
	"time"
)

type orgUsecase struct {
	repo repositories.OrgRepository
}

type OrgUsecase interface {
	GetOrgHash(ctx context.Context, req clDto.GetOrgHashRequest) (map[string]any, error)
	GetOrgData(ctx context.Context, req clDto.GetOrgDataRequest) (bool, interface{}, error)

	ServerCreateDept(ctx context.Context, req svDto.SvCreateDeptRequest) (interface{}, error)
	ServerDeleteDept(ctx context.Context, req svDto.SvDeleteDeptRequest) (interface{}, error)

	ServerCreateOrgFile(ctx context.Context, req svDto.SvCreateOrgFileRequest) (interface{}, error)

	ServerCreateDeptUser(ctx context.Context, req svDto.SvCreateDeptUserRequest) (interface{}, error)
	ServerDeleteDeptUser(ctx context.Context, req svDto.SvDeleteDeptUserRequest) (interface{}, error)
}

func NewOrgUsecase(repo repositories.OrgRepository) OrgUsecase {
	return &orgUsecase{repo: repo}
}

func (r *orgUsecase) GetOrgHash(ctx context.Context, req clDto.GetOrgHashRequest) (map[string]any, error) {

	orgMap := make(map[string]any)

	for i := 0; i < len(req.OrgHash); i++ {
		parts := strings.Split(req.OrgHash[i], "_")
		if len(parts) == 2 {

			fileFlag, eventFlag, err := r.repo.CheckOrgHash(ctx, parts[0], parts[1])

			if err != nil {
				return nil, err
			} else if fileFlag {
				// 파일로 받아야함.
				orgMap[req.OrgHash[i]] = "file"
			} else if eventFlag {
				// 이벤트로 받아야함.
				orgMap[req.OrgHash[i]] = "event"
			} else {
				// 최신 버전.
				orgMap[req.OrgHash[i]] = "latest"
			}
		} else {
			fmt.Printf("GetOrgs org : %s is invalid !", req.OrgHash[i])
			continue
		}
	}

	return orgMap, nil
}

// func parseToOrgEventEntities(models []models.OrgEvent) []entities.OrgEventEntity {
// 	var eventList []entities.OrgEventEntity

// 	for _, m := range models {
// 		entity := entities.OrgEventEntity{
// 			Seq:        m.Seq,
// 			OrgCode:    m.OrgCode,
// 			Kind:       m.Kind,
// 			EventType:  m.EventType,
// 			UpdateHash: m.UpdateHash,
// 		}
// 		eventList = append(eventList, entity)
// 	}
// 	return eventList
// }

// 클라이언트, 서버에서 요청하여 현재의 ORG를 가져오기 위해 각자의 타입에 맞춰 entity 생성하는 코드 (오버라이드를 지원하지 않기 때문에 사용함)
func (r *orgUsecase) toGetOrgEntity(orgCode string) entities.GetOrgEntity {
	return entities.GetOrgEntity{
		OrgCode: orgCode,
	}
}

func parseOrgTree(orgTree []models.WorksOrg) *entities.OrgEntity {

	if orgTree == nil {
		fmt.Println("조회된 조직도 정보가 없음. ")
		return nil
	}

	// 최상위 구조
	var rootOrgInfos []entities.OrgInfo
	var flatList []entities.OrgInfo // 트리 구성용 전체 flat 리스트

	for _, org := range orgTree {
		// 이름 다국어 처리
		name := entities.NameEntity{
			Kr: org.KrLang,
			En: org.EnLang,
			Jp: org.JpLang,
			Cn: org.CnLang,
		}

		info := entities.OrgInfo{
			DeptCode:       org.DeptCode,
			ParentDeptCode: org.ParentDeptCode,
			Name:           name,
			Kind:           org.Kind,
		}

		if org.ParentDeptCode == "root" {
			rootOrgInfos = append(rootOrgInfos, info)
		}
		flatList = append(flatList, info)
	}

	// 트리 구조로 변환
	orgTreeInfos := buildOrgTree(flatList, "root")

	return &entities.OrgEntity{
		RootDept: rootOrgInfos,
		OrgTree:  orgTreeInfos,
	}
}

func buildOrgTree(flatList []entities.OrgInfo, parentCode string) []entities.OrgTreeInfos {
	var tree []entities.OrgTreeInfos

	for _, org := range flatList {
		if org.ParentDeptCode == parentCode {
			// 재귀적으로 하위 부서를 구성
			sub := buildOrgTree(flatList, org.DeptCode)

			tree = append(tree, entities.OrgTreeInfos{
				DeptCode:       org.DeptCode,
				ParentDeptCode: org.ParentDeptCode,
				Name:           org.Name,
				SubDept:        sub,
				Kind:           org.Kind,
			})
		}
	}

	return tree
}

func (r *orgUsecase) ServerCreateDept(ctx context.Context, req svDto.SvCreateDeptRequest) (interface{}, error) {
	return r.repo.SaveDept(ctx, toCreateDepartmentEntity(req))
}

func toCreateDepartmentEntity(req svDto.SvCreateDeptRequest) entities.CreateDeptEntity {

	return entities.CreateDeptEntity{
		DeptCode:       req.DeptCode,
		DeptOrg:        req.DeptOrg,
		ParentDeptCode: req.ParentDeptCode,
		KrLang:         req.KrLang,
		EnLang:         req.EnLang,
		JpLang:         req.JpLang,
		CnLang:         req.CnLang,
	}
}

func (r *orgUsecase) ServerDeleteDept(ctx context.Context, req svDto.SvDeleteDeptRequest) (interface{}, error) {
	return r.repo.DeleteDept(ctx, toDeleteDepartmentEntity(req))
}

func toDeleteDepartmentEntity(req svDto.SvDeleteDeptRequest) entities.DeleteDeptEntity {

	return entities.DeleteDeptEntity{
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}

func (r *orgUsecase) ServerCreateOrgFile(ctx context.Context, req svDto.SvCreateOrgFileRequest) (interface{}, error) {

	for i := 0; i < len(req.OrgCode); i++ {

		orgTree, err := r.repo.GetOrg(ctx, r.toGetOrgEntity(req.OrgCode[i]))

		if err != nil {
			fmt.Printf("ServerCreateOrgFile org : %s is invalid ! \n", req.OrgCode[i])
			continue
		}

		// 파일 명 생성.
		fileName := req.OrgCode[i] + "_" + getNow()
		fmt.Printf("org %s file name : %s ", req.OrgCode[i], fileName)

		orgEntity := parseOrgTree(orgTree)
		// orgEntity를 JSON 등으로 직렬화 하고 내용을 ZIP 파일 내에 저장

		// 1. OrgEntity → JSON 직렬화
		fmt.Println("OrgEntity → JSON 직렬화")
		orgJson, err := json.MarshalIndent(orgEntity, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal OrgEntity: %w", err)
		}
		fmt.Println("OrgEntity → JSON 직렬화 ok")

		// 경로에 파일명을 포함시켜야 함
		var zipPath = "./storage/org_files/" + fileName

		// 1. 디렉터리 없으면 생성 (디렉터리 경로만 던져야함.) -> 실행시 마운트 필요
		err = ensureDir("./storage/org_files/")
		if err != nil {
			return nil, fmt.Errorf("디렉터리 생성 실패: %w", err)
		}

		// 2. ZIP 파일 생성
		fmt.Println("ZIP 파일 생성")
		zipFile, err := os.Create(zipPath)
		if err != nil {
			fmt.Println(err.Error())
			return nil, fmt.Errorf("failed to create zip file: %w", err)
		}
		defer zipFile.Close()
		fmt.Println("ZIP 파일 생성 ok")

		// 3. ZIP writer 생성
		fmt.Println("ZIP writer 생성")
		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()
		fmt.Println("ZIP writer 생성 ok")

		// 4. ZIP 내 파일 생성
		fmt.Println("ZIP 내 파일 생성")
		fileWriter, err := zipWriter.Create(fileName)
		if err != nil {
			return nil, fmt.Errorf("failed to create file in zip: %w", err)
		}
		fmt.Println("ZIP 내 파일 생성 ok")

		// 5. write
		fmt.Println("Write")
		_, err = fileWriter.Write(orgJson)
		if err != nil {
			return nil, fmt.Errorf("failed to write data to zip: %w", err)
		}
		fmt.Println("Write ok")

	}
	return consts.SUCCESS, nil
}

func getNow() string {
	now := time.Now()
	formatted := now.Format(consts.YYYYMMDDHHMSS)
	return formatted
}

// 경로가 없으면 생성
func ensureDir(dirPath string) error {
	return os.MkdirAll(dirPath, os.ModePerm)
}

func (r *orgUsecase) GetOrgData(ctx context.Context, req clDto.GetOrgDataRequest) (bool, interface{}, error) {

	if req.Type == consts.FILE {
		version, err := r.repo.GetOrgLatestVersion(ctx, req.OrgCode)
		if err != nil {
			return false, nil, err
		}

		filePath := "./storage/" + req.OrgCode + "/org_files/" + version // 전달할 파일 경로
		// 파일을 메모리에 가지고 있도록 수정 할 것.
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("파일을 찾을 수 없음 %s \n", filePath)
			return false, nil, err
		}
		return true, fileBytes, nil

	} else if req.Type == consts.EVENT {

		events, err := r.repo.GetOrgDiffEvent(ctx, req.OrgCode, req.OrgHash)
		if err != nil {
			return false, nil, err
		}
		return false, events, nil

	} else {
		// 명확하지 않은 타입으로 요청함.
		return false, nil, fmt.Errorf("invalid request type")
	}

}

func (r *orgUsecase) ServerCreateDeptUser(ctx context.Context, req svDto.SvCreateDeptUserRequest) (interface{}, error) {

	updateHash := makeUpdateHash()
	fmt.Println("사용자 추가시 update Hash 생성 : ", updateHash)

	return r.repo.SaveDeptUser(ctx, toCreateDeptUserEntity(req, updateHash))
}

func toCreateDeptUserEntity(req svDto.SvCreateDeptUserRequest, updateHash string) entities.CreateDeptUserEntity {

	return entities.CreateDeptUserEntity{
		UserHash:             req.UserHash,
		DeptCode:             req.DeptCode,
		DeptOrg:              req.DeptOrg,
		PositionCode:         req.PositionCode,
		RoleCode:             req.RoleCode,
		IsConcurrentPosition: req.IsConcurrentPosition,
		UpdateHash:           updateHash,
	}
}

func makeUpdateHash() string {
	// 현재 시간 밀리초 문자열
	now := time.Now().UnixNano() // 나노초 단위 시간값

	// 16바이트 랜덤 바이트 생성
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	// 시간값을 바이트 배열로 변환 (int64 -> []byte)
	timeBytes := []byte(fmt.Sprintf("%d", now))

	// 시간 + 랜덤 바이트 합치기
	data := append(timeBytes, randomBytes...)

	// SHA-256 해시 생성
	hash := sha256.Sum256(data)

	// 16진수 인코딩해서 문자열 반환
	return hex.EncodeToString(hash[:])
}

func (r *orgUsecase) ServerDeleteDeptUser(ctx context.Context, req svDto.SvDeleteDeptUserRequest) (interface{}, error) {
	return r.repo.DeleteDeptUser(ctx, toDeleteDeptUserEntity(req))
}

func toDeleteDeptUserEntity(req svDto.SvDeleteDeptUserRequest) entities.DeleteDeptUserEntity {

	return entities.DeleteDeptUserEntity{
		UserHash: req.UserHash,
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}
