package usecase

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"org/internal/application/usecase/input"
	"org/internal/application/util"
	"org/internal/consts"
	"org/internal/delivery/dto/org"
	"org/internal/domain/org/entity"
	"org/internal/domain/org/repository"
	sharedEntity "org/internal/domain/shared/entity"
	"org/internal/infrastructure/storage"
	commonConsts "org/pkg/consts"
	"strings"
)

type orgUsecase struct {
	repository     repository.OrgRepository
	orgFileStorage storage.OrgFileStorage
	orgStorage     storage.OrgStorage
}

type OrgUsecase interface {
	GetOrgHash(ctx context.Context, req org.GetOrgHashRequest) (map[string]any, error)
	GetOrgData(ctx context.Context, req org.GetOrgDataRequest) (string, interface{}, error)
	CreateOrgFile(ctx context.Context, req org.CreateOrgFileRequest) (interface{}, error)
	RegistOrgBatch(ctx context.Context, in input.RegistOrgBatchInput) error
	GetWorksOrgCode() []string
}

func NewOrgUsecase(repository repository.OrgRepository, orgFileStorage storage.OrgFileStorage, orgStorage storage.OrgStorage) OrgUsecase {
	return &orgUsecase{
		repository:     repository,
		orgFileStorage: orgFileStorage,
		orgStorage:     orgStorage,
	}
}

func (r *orgUsecase) GetOrgHash(ctx context.Context, req org.GetOrgHashRequest) (map[string]any, error) {

	orgMap := make(map[string]any)

	for i := 0; i < len(req.OrgHash); i++ {
		parts := strings.Split(req.OrgHash[i], "_")
		if len(parts) == 2 {

			fileFlag, eventFlag, err := r.repository.CheckOrgHash(ctx, parts[0], parts[1])

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

func (r *orgUsecase) GetOrgData(ctx context.Context, req org.GetOrgDataRequest) (string, interface{}, error) {

	serverOrgCode := r.orgStorage.WorksOrgCodeExist(req.OrgCode)

	log.Println("serverOrgCode : ", serverOrgCode)
	if !serverOrgCode {
		log.Printf("[GetOrgData] orgCode %s not exist. \n", req.OrgCode)
		return "", nil, consts.ErrOrgCodeNotExist
	}

	if req.Type == consts.FILE {
		version, err := r.repository.GetOrgLatestVersion(ctx, req.OrgCode)
		if err != nil {
			return "", nil, err
		}

		data, err := r.orgFileStorage.GetOrgFile(req.OrgCode)

		//filePath := "./storage/" + req.OrgCode + "/org_files/" + version // 전달할 파일 경로
		// 파일을 메모리에 가지고 있도록 수정 할 것.
		// fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("파일을 찾을 수 없음 %s \n", req.OrgCode)
			return "", nil, consts.ErrOrgFileNotFound
		}
		return version, data, nil

	} else if req.Type == consts.EVENT {

		events, err := r.repository.GetOrgDiffEvent(ctx, req.OrgCode, req.OrgHash)
		if err != nil {
			return "", nil, err
		}
		return "", events, nil

	} else {
		// 명확하지 않은 타입으로 요청함.
		return "", nil, fmt.Errorf("invalid request type")
	}

}

func (r *orgUsecase) CreateOrgFile(ctx context.Context, req org.CreateOrgFileRequest) (interface{}, error) {

	for i := 0; i < len(req.OrgCode); i++ {

		org := req.OrgCode[i]

		orgTree, err := r.repository.GetOrg(ctx, org)
		if err != nil {
			fmt.Printf("Invalid org: %s\n", req.OrgCode[i])
			continue
		}

		// 저장시간 생성 = 파일 명
		fileName := util.GetNow() + ".json"
		fmt.Printf("org %s file name: %s\n", org, fileName)

		orgEntity := parseOrgTree(orgTree)
		orgJson, err := json.MarshalIndent(orgEntity, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshal error: %w", err)
		}

		// 메모리에서 ZIP 생성
		zipData, err := buildZipInMemory(fileName, orgJson)
		if err != nil {
			return nil, fmt.Errorf("zip build error: %w", err)
		}

		// 메모리 저장소에 저장
		if err := r.orgFileStorage.SaveOrgFile(org, zipData); err != nil {
			return nil, fmt.Errorf("memory save error: %w", err)
		}

		// 점검
		data, err := r.orgFileStorage.GetOrgFile(org)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("저장된 org 파일의 사이즈 :", len(data))

		// DB 저장
		if ok, err := r.repository.PutOrgEventHash(ctx, org, fileName); err != nil {
			return nil, fmt.Errorf("db save error: %w", err)
		} else if ok {
			log.Println("DB saved ok org:", org)
		}
	}

	return commonConsts.SUCCESS, nil
}

func buildZipInMemory(fileName string, content []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	writer, err := zipWriter.Create(fileName)
	if err != nil {
		return nil, err
	}
	if _, err := writer.Write(content); err != nil {
		return nil, err
	}
	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func parseOrgTree(orgTree []entity.WorksOrg) *entity.OrgEntity {

	if orgTree == nil {
		log.Println("조회된 조직도 정보가 없음. ")
		return nil
	}

	// 최상위 구조
	var rootOrgInfos []entity.OrgInfo
	var flatList []entity.OrgInfo // 트리 구성용 전체 flat 리스트

	for _, org := range orgTree {
		// 이름 다국어 처리
		name := sharedEntity.NameEntity{
			Def: org.KoLang, // 수정 필요.
			Ko:  org.KoLang,
			En:  org.EnLang,
			Jp:  org.JpLang,
			Zh:  org.ZhLang,
			Ru:  org.RuLang,
			Vi:  org.ViLang,
		}

		info := entity.OrgInfo{
			DeptCode:       org.DeptCode,
			ParentDeptCode: org.ParentDeptCode,
			Name:           name,
			Kind:           org.Kind,
			UserHash:       org.UserHash,
			UserId:         org.UserId,
			Header:         org.Header,
			Description:    org.Description,
		}

		if org.ParentDeptCode == "root" {
			rootOrgInfos = append(rootOrgInfos, info)
		}
		flatList = append(flatList, info)
	}

	// 트리 구조로 변환
	orgTreeInfo := buildOrgTree(flatList, "root")

	return &entity.OrgEntity{
		RootDept: rootOrgInfos,
		OrgTree:  orgTreeInfo,
	}
}

func buildOrgTree(flatList []entity.OrgInfo, parentCode string) []entity.OrgTreeInfo {
	var tree []entity.OrgTreeInfo

	for _, org := range flatList {
		if org.ParentDeptCode == parentCode {
			// 재귀적으로 하위 부서를 구성
			sub := buildOrgTree(flatList, org.DeptCode)

			// 사실 이렇게 구분해서 init하지 않아도 entity.OrgTreeInfo 내부에서 omitempty처리하면 response시 보이지 않음.
			if org.Kind == "0" {
				// 부서
				tree = append(tree, entity.OrgTreeInfo{
					DeptCode:       org.DeptCode,
					ParentDeptCode: org.ParentDeptCode,
					Name:           org.Name,
					SubDept:        sub,
					Kind:           org.Kind,
					UserHash:       org.UserHash,
					Description:    org.Description,
				})

				// 사용자
			} else if org.Kind == "1" {
				tree = append(tree, entity.OrgTreeInfo{
					ParentDeptCode: org.ParentDeptCode,
					Name:           org.Name,
					SubDept:        sub,
					Kind:           org.Kind,
					UserHash:       org.UserHash,
					UserId:         org.UserId,
					Header:         org.Header,
				})
			}

		}
	}

	return tree
}

func (r *orgUsecase) RegistOrgBatch(ctx context.Context, in input.RegistOrgBatchInput) error {

	en := entity.MakeRegistOrgBatchEntity(in.File, in.FileName, in.OrgCode)

	// zip 해제, json 구하기
	jsonBytes, err := unzipAndGetJSON(en.OrgFile)

	if err != nil {
		log.Println("[RegistOrgBatch] unzipAndGetJSON error")
		return consts.ErrUnzipAndGetJSONError
	}

	// 2. JSON → Wrapper
	var orgInfo []entity.WorksOrg
	if err := json.Unmarshal(jsonBytes, &orgInfo); err != nil {
		return consts.ErrInvalidOrgJSONError
	}

	// 3. diff 구하기 (현재 데이터 조회)
	// nowOrgInfo, err := r.repository.GetOrg(ctx, en.OrgCode)
	// if err != nil {
	// 	return err
	// }

	// 부서 - 하위 사용자 구조로 만들고, 부서 * 부서 하위사용자 수만큼 반복

	dept, user := splitByKind(orgInfo)

	// DB 저장 insert update.
	err = r.repository.RegistOrgBatch(ctx, dept, user)
	if err != nil {
		return err
	}

	return nil
}

func unzipAndGetJSON(orgFile *[]byte) ([]byte, error) {

	// zip reader 생성
	zr, err := zip.NewReader(
		bytes.NewReader(*orgFile),
		int64(len(*orgFile)),
	)
	if err != nil {
		return nil, fmt.Errorf("zip reader error: %w", err)
	}

	// ZIP 내부 파일 순회
	for _, f := range zr.File {

		// json 파일만 추출
		if !strings.HasSuffix(f.Name, ".json") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()

		// json 읽기
		jsonBytes, err := io.ReadAll(rc)
		if err != nil {
			return nil, err
		}

		return jsonBytes, nil
	}

	return nil, fmt.Errorf("json file not found in zip")
}

func splitByKind(temp []entity.WorksOrg) ([]entity.WorksOrg, []entity.WorksOrg) {

	depts := make([]entity.WorksOrg, 0)
	users := make([]entity.WorksOrg, 0)

	for _, item := range temp {
		switch item.Kind {
		case "0": // 부서
			depts = append(depts, entity.WorksOrg{
				Org:            item.Org,
				DeptCode:       item.DeptCode,
				ParentDeptCode: item.ParentDeptCode,
				KoLang:         item.KoLang,
				EnLang:         item.EnLang,
				UpdateHash:     item.UpdateHash,
				Header:         item.Header,
				Description:    item.Description,
			})

		case "1": // 사용자
			users = append(users, entity.WorksOrg{
				Org:        item.Org,
				UserHash:   item.UserHash,
				UserId:     item.UserId,
				DeptCode:   item.ParentDeptCode, // 사용자는 상위 부서
				KoLang:     item.KoLang,
				EnLang:     item.EnLang,
				UpdateHash: item.UpdateHash,
			})
		}
	}

	return depts, users
}

func (r *orgUsecase) GetWorksOrgCode() []string {
	return r.orgStorage.GetWorksOrgCode()
}
