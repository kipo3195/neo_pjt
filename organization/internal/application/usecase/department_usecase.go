package usecase

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"org/internal/application/usecase/input"
	"org/internal/application/util"
	"org/internal/delivery/dto/department"
	"org/internal/domain/department/entity"
	"org/internal/domain/department/repository"
)

type departmentUsecase struct {
	repository repository.DepartmentRepository
}

type DepartmentUsecase interface {
	CreateDept(ctx context.Context, input input.CreateDeptInput) (interface{}, error)

	DeleteDept(ctx context.Context, req department.DeleteDeptRequest) (interface{}, error)

	CreateDeptUser(ctx context.Context, req department.CreateDeptUserRequest) (interface{}, error)
	DeleteDeptUser(ctx context.Context, req department.DeleteDeptUserRequest) (interface{}, error)

	CreateWorksDept(ctx context.Context, input input.CreateWorksDeptInput) error

	CreateWorksDeptMultiLang(ctx context.Context, input input.CreateWorksDeptMultiLangInput) error
	CreateDeptMultiLang(ctx context.Context, en entity.WorksDeptEntity) error

	CreateWorksDeptUser(ctx context.Context, input input.CreateWorksDeptUserInput) error
}

func NewDepartmentUsecase(repository repository.DepartmentRepository) DepartmentUsecase {
	return &departmentUsecase{
		repository: repository,
	}
}

func (r *departmentUsecase) CreateDept(ctx context.Context, in input.CreateDeptInput) (interface{}, error) {

	entity := entity.MakeCreateDeptEntity(in.DeptCode, in.DeptOrg, in.ParentDeptCode, in.KoLang, in.EnLang, in.JpLang, in.RuLang, in.ViLang, in.ZhLang, in.Header)
	return r.repository.PutDept(ctx, entity)
}

func (r *departmentUsecase) DeleteDeptUser(ctx context.Context, req department.DeleteDeptUserRequest) (interface{}, error) {
	return r.repository.DeleteDeptUser(ctx, toDeleteDeptUserEntity(req))
}

func toDeleteDeptUserEntity(req department.DeleteDeptUserRequest) entity.DeleteDeptUserEntity {

	return entity.DeleteDeptUserEntity{
		UserHash: req.UserHash,
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}

func (r *departmentUsecase) CreateDeptUser(ctx context.Context, req department.CreateDeptUserRequest) (interface{}, error) {

	updateHash := util.MakeUpdateHash()
	log.Println("사용자 추가시 update Hash 생성 : ", updateHash)

	return r.repository.PutDeptUser(ctx, toCreateDeptUserEntity(req, updateHash))
}

func toCreateDeptUserEntity(req department.CreateDeptUserRequest, updateHash string) entity.CreateDeptUserEntity {

	return entity.CreateDeptUserEntity{
		UserHash:             req.UserHash,
		DeptCode:             req.DeptCode,
		DeptOrg:              req.DeptOrg,
		PositionCode:         req.PositionCode,
		RoleCode:             req.RoleCode,
		IsConcurrentPosition: req.IsConcurrentPosition,
		UpdateHash:           updateHash,
	}
}

func (r *departmentUsecase) DeleteDept(ctx context.Context, req department.DeleteDeptRequest) (interface{}, error) {
	return r.repository.DeleteDept(ctx, toDeleteDepartmentEntity(req))
}

func toDeleteDepartmentEntity(req department.DeleteDeptRequest) entity.DeleteDeptEntity {

	return entity.DeleteDeptEntity{
		DeptCode: req.DeptCode,
		DeptOrg:  req.DeptOrg,
	}
}

func (u *departmentUsecase) CreateWorksDept(ctx context.Context, input input.CreateWorksDeptInput) error {

	u.CreateDeptTree(ctx, "root", 1, input.MaxDepth, input.Org, input.DeptCount)

	return nil

}

func (u *departmentUsecase) CreateDeptTree(ctx context.Context, parentCode string, depth, maxDepth int, org string, deptCount int) error {
	if depth > maxDepth {
		return nil
	}

	numDepts := rand.Intn(deptCount) + 1 // 각 단계에서 1~5개의 부서
	for i := 0; i < numDepts; i++ {
		deptCode := randomString(8)
		log.Printf("depth : %d, code : %s", depth, deptCode)
		dept := entity.WorksDeptEntity{
			DeptCode:        deptCode,
			DeptOrg:         org,
			ParentsDeptCode: parentCode,
		}

		if err := u.repository.CreateDeptTree(ctx, dept); err != nil {
			return fmt.Errorf("insert failed (deptCode=%s): %w", deptCode, err)
		}

		// 하위 부서 생성
		if err := u.CreateDeptTree(ctx, deptCode, depth+1, maxDepth, org, deptCount); err != nil {
			return err
		}
	}
	return nil
}

func randomString(n int) string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func (u *departmentUsecase) CreateWorksDeptMultiLang(ctx context.Context, input input.CreateWorksDeptMultiLangInput) error {

	en, err := u.repository.GetDepts(ctx, input.Org)

	if err != nil {
		return err
	}

	for i := 0; i < len(en); i++ {
		u.CreateDeptMultiLang(ctx, en[i])
	}

	return nil
}

func (u *departmentUsecase) CreateDeptMultiLang(ctx context.Context, en entity.WorksDeptEntity) error {

	// 🔹 한글
	koreanParts := []string{"개발팀", "기획팀", "디자인팀", "마케팅팀", "인사팀", "재무팀", "영업팀", "고객지원팀", "보안팀", "품질관리팀",
		"데이터분석팀", "AI연구팀", "클라우드팀", "서버운영팀", "네트워크팀", "기술지원팀", "법무팀", "홍보팀", "교육팀", "전략기획팀",
		"UX팀", "UI팀", "모바일개발팀", "웹개발팀", "백엔드팀", "프론트엔드팀", "인프라팀", "QA팀", "테스트팀", "프로덕트팀",
		"운영팀", "구매팀", "생산관리팀", "품질보증팀", "R&D팀", "연구소", "해외영업팀", "국내영업팀", "CS팀", "서비스기획팀",
		"전산팀", "기술연구소", "기획운영팀", "콘텐츠팀", "브랜드팀", "디지털전략팀", "광고팀", "SNS운영팀", "파트너관리팀", "협력사관리팀",
		"출판팀", "영상제작팀", "촬영팀", "사운드팀", "이벤트기획팀", "행사운영팀", "총무팀", "시설관리팀", "보안운영팀", "데이터엔지니어팀",
		"데이터사이언스팀", "AI모델팀", "ML팀", "로보틱스팀", "IoT팀", "UI/UX리서치팀", "서비스디자인팀", "제품디자인팀", "콘텐츠기획팀", "CRM팀",
		"광고운영팀", "마케팅전략팀", "퍼포먼스마케팅팀", "리서치팀", "인재개발팀", "채용팀", "보상팀", "노무팀", "교육운영팀", "복지팀",
		"회계팀", "세무팀", "예산팀", "자금팀", "투자팀", "IR팀", "법무지원팀", "리스크관리팀", "내부감사팀", "감사팀",
		"전략혁신팀", "기술혁신팀", "신사업팀", "신제품팀", "서비스운영팀", "클라이언트지원팀", "파트너기술팀", "모바일서비스팀", "웹서비스팀", "플랫폼팀",
		"API개발팀", "데브옵스팀", "클라우드인프라팀", "보안개발팀", "암호화팀", "인증팀", "DB운영팀", "데이터관리팀", "AI시스템팀", "모델서빙팀",
		"로봇개발팀", "시스템설계팀", "하드웨어팀", "펌웨어팀", "품질보증연구소", "테크랩", "AI솔루션팀", "프로젝트관리팀", "PMO팀", "사업기획팀",
		"신규사업팀", "파트너개발팀", "서비스디자인랩", "글로벌전략팀", "해외사업팀", "수출팀", "해외지원팀", "통번역팀", "문서관리팀", "지식관리팀",
		"IT운영팀", "헬프데스크팀", "보안정책팀", "인증운영팀", "접근관리팀", "네트워크보안팀", "사이버보안팀", "데이터보호팀", "프라이버시팀", "리스크평가팀",
		"고객경험팀", "CX팀", "사용자리서치팀", "VOC팀", "고객만족팀", "커뮤니티운영팀", "브랜드전략팀", "브랜드디자인팀", "프로모션팀", "캠페인팀",
		"온라인마케팅팀", "오프라인마케팅팀", "디지털콘텐츠팀", "영상콘텐츠팀", "SNS콘텐츠팀", "뉴스레터팀", "광고디자인팀", "일러스트팀", "3D디자인팀", "프린트디자인팀",
		"사내홍보팀", "대외협력팀", "정부협력팀", "정책기획팀", "산학협력팀", "연구기획팀", "특허팀", "인증획득팀", "기술평가팀", "품질인증팀",
		"지속가능경영팀", "ESG팀", "환경안전팀", "산업안전팀", "에너지관리팀", "시설안전팀", "보건팀", "응급대응팀", "출입관리팀", "보안감시팀",
		"문서보안팀", "개인정보보호팀", "정보관리팀", "지식공유팀", "학습조직팀", "리더십개발팀", "성과관리팀", "성과보상팀", "교육개발팀", "HR전략팀"}

	// 🔹 영어 이름 파츠
	englishParts := []string{"Development Team", "Planning Team", "Design Team", "Marketing Team", "HR Team", "Finance Team", "Sales Team", "Customer Support Team", "Security Team", "Quality Control Team",
		"Data Analysis Team", "AI Research Team", "Cloud Team", "Server Operations Team", "Network Team", "Technical Support Team", "Legal Team", "PR Team", "Training Team", "Strategic Planning Team",
		"UX Team", "UI Team", "Mobile Development Team", "Web Development Team", "Backend Team", "Frontend Team", "Infrastructure Team", "QA Team", "Testing Team", "Product Team",
		"Operations Team", "Procurement Team", "Production Management Team", "Quality Assurance Team", "R&D Team", "Research Center", "Overseas Sales Team", "Domestic Sales Team", "Customer Service Team", "Service Planning Team",
		"IT Operations Team", "Technical Research Center", "Planning and Operations Team", "Content Team", "Brand Team", "Digital Strategy Team", "Advertising Team", "Social Media Management Team", "Partner Management Team", "Vendor Management Team",
		"Publishing Team", "Video Production Team", "Filming Team", "Sound Team", "Event Planning Team", "Event Operations Team", "General Affairs Team", "Facility Management Team", "Security Operations Team", "Data Engineering Team",
		"Data Science Team", "AI Model Team", "Machine Learning Team", "Robotics Team", "IoT Team", "UI/UX Research Team", "Service Design Team", "Product Design Team", "Content Planning Team", "CRM Team",
		"Ad Operations Team", "Marketing Strategy Team", "Performance Marketing Team", "Research Team", "Talent Development Team", "Recruitment Team", "Compensation Team", "Labor Relations Team", "Training Operations Team", "Welfare Team",
		"Accounting Team", "Tax Team", "Budget Team", "Fund Management Team", "Investment Team", "IR Team", "Legal Support Team", "Risk Management Team", "Internal Audit Team", "Audit Team",
		"Strategic Innovation Team", "Technical Innovation Team", "New Business Team", "New Product Team", "Service Operations Team", "Client Support Team", "Partner Technology Team", "Mobile Service Team", "Web Service Team", "Platform Team",
		"API Development Team", "DevOps Team", "Cloud Infrastructure Team", "Security Development Team", "Encryption Team", "Authentication Team", "Database Operations Team", "Data Management Team", "AI Systems Team", "Model Serving Team",
		"Robotics Development Team", "System Design Team", "Hardware Team", "Firmware Team", "Quality Assurance Lab", "Tech Lab", "AI Solutions Team", "Project Management Team", "PMO Team", "Business Planning Team",
		"New Project Team", "Partner Development Team", "Service Design Lab", "Global Strategy Team", "Overseas Business Team", "Export Team", "International Support Team", "Translation and Interpretation Team", "Document Management Team", "Knowledge Management Team",
		"IT Operations Team", "Helpdesk Team", "Security Policy Team", "Authentication Operations Team", "Access Management Team", "Network Security Team", "Cyber Security Team", "Data Protection Team", "Privacy Team", "Risk Assessment Team",
		"Customer Experience Team", "CX Team", "User Research Team", "Voice of Customer Team", "Customer Satisfaction Team", "Community Operations Team", "Brand Strategy Team", "Brand Design Team", "Promotion Team", "Campaign Team",
		"Online Marketing Team", "Offline Marketing Team", "Digital Content Team", "Video Content Team", "SNS Content Team", "Newsletter Team", "Advertising Design Team", "Illustration Team", "3D Design Team", "Print Design Team",
		"Internal PR Team", "External Relations Team", "Government Relations Team", "Policy Planning Team", "Industry-Academia Cooperation Team", "Research Planning Team", "Patent Team", "Certification Acquisition Team", "Technical Evaluation Team", "Quality Certification Team",
		"Sustainability Management Team", "ESG Team", "Environment and Safety Team", "Industrial Safety Team", "Energy Management Team", "Facility Safety Team", "Health Team", "Emergency Response Team", "Access Control Team", "Security Monitoring Team",
		"Document Security Team", "Personal Information Protection Team", "Information Management Team", "Knowledge Sharing Team", "Learning Organization Team", "Leadership Development Team", "Performance Management Team", "Compensation and Rewards Team", "Training Development Team", "HR Strategy Team"}

	// 🔹 베트남 이름 파츠
	vietnameseParts := []string{"Nhóm Phát triển", "Nhóm Kế hoạch", "Nhóm Thiết kế", "Nhóm Tiếp thị", "Nhóm Nhân sự", "Nhóm Tài chính", "Nhóm Bán hàng", "Nhóm Hỗ trợ khách hàng", "Nhóm Bảo mật", "Nhóm Kiểm soát chất lượng",
		"Nhóm Phân tích dữ liệu", "Nhóm Nghiên cứu AI", "Nhóm Đám mây", "Nhóm Vận hành máy chủ", "Nhóm Mạng", "Nhóm Hỗ trợ kỹ thuật", "Nhóm Pháp lý", "Nhóm Quan hệ công chúng", "Nhóm Đào tạo", "Nhóm Lập kế hoạch chiến lược",
		"Nhóm UX", "Nhóm UI", "Nhóm Phát triển di động", "Nhóm Phát triển web", "Nhóm Backend", "Nhóm Frontend", "Nhóm Cơ sở hạ tầng", "Nhóm QA", "Nhóm Kiểm thử", "Nhóm Sản phẩm",
		"Nhóm Vận hành", "Nhóm Mua hàng", "Nhóm Quản lý sản xuất", "Nhóm Đảm bảo chất lượng", "Nhóm R&D", "Trung tâm nghiên cứu", "Nhóm Kinh doanh quốc tế", "Nhóm Kinh doanh nội địa", "Nhóm CS", "Nhóm Lập kế hoạch dịch vụ",
		"Nhóm CNTT", "Phòng Nghiên cứu kỹ thuật", "Nhóm Kế hoạch và vận hành", "Nhóm Nội dung", "Nhóm Thương hiệu", "Nhóm Chiến lược kỹ thuật số", "Nhóm Quảng cáo", "Nhóm Quản lý mạng xã hội", "Nhóm Quản lý đối tác", "Nhóm Quản lý nhà cung cấp",
		"Nhóm Xuất bản", "Nhóm Sản xuất video", "Nhóm Quay phim", "Nhóm Âm thanh", "Nhóm Lên kế hoạch sự kiện", "Nhóm Vận hành sự kiện", "Nhóm Hành chính tổng hợp", "Nhóm Quản lý cơ sở vật chất", "Nhóm Vận hành bảo mật", "Nhóm Kỹ sư dữ liệu",
		"Nhóm Khoa học dữ liệu", "Nhóm Mô hình AI", "Nhóm Học máy", "Nhóm Robot", "Nhóm IoT", "Nhóm Nghiên cứu UI/UX", "Nhóm Thiết kế dịch vụ", "Nhóm Thiết kế sản phẩm", "Nhóm Lên kế hoạch nội dung", "Nhóm CRM",
		"Nhóm Quản lý quảng cáo", "Nhóm Chiến lược tiếp thị", "Nhóm Tiếp thị hiệu suất", "Nhóm Nghiên cứu", "Nhóm Phát triển nhân tài", "Nhóm Tuyển dụng", "Nhóm Lương thưởng", "Nhóm Lao động", "Nhóm Đào tạo", "Nhóm Phúc lợi",
		"Nhóm Kế toán", "Nhóm Thuế", "Nhóm Ngân sách", "Nhóm Quản lý quỹ", "Nhóm Đầu tư", "Nhóm IR", "Nhóm Hỗ trợ pháp lý", "Nhóm Quản lý rủi ro", "Nhóm Kiểm toán nội bộ", "Nhóm Kiểm toán",
		"Nhóm Đổi mới chiến lược", "Nhóm Đổi mới công nghệ", "Nhóm Kinh doanh mới", "Nhóm Sản phẩm mới", "Nhóm Vận hành dịch vụ", "Nhóm Hỗ trợ khách hàng", "Nhóm Kỹ thuật đối tác", "Nhóm Dịch vụ di động", "Nhóm Dịch vụ web", "Nhóm Nền tảng",
		"Nhóm Phát triển API", "Nhóm DevOps", "Nhóm Hạ tầng đám mây", "Nhóm Phát triển bảo mật", "Nhóm Mã hóa", "Nhóm Xác thực", "Nhóm Vận hành cơ sở dữ liệu", "Nhóm Quản lý dữ liệu", "Nhóm Hệ thống AI", "Nhóm Triển khai mô hình",
		"Nhóm Phát triển robot", "Nhóm Thiết kế hệ thống", "Nhóm Phần cứng", "Nhóm Phần mềm nhúng", "Phòng Đảm bảo chất lượng", "Phòng Công nghệ", "Nhóm Giải pháp AI", "Nhóm Quản lý dự án", "Nhóm PMO", "Nhóm Kế hoạch kinh doanh",
		"Nhóm Dự án mới", "Nhóm Phát triển đối tác", "Phòng Thiết kế dịch vụ", "Nhóm Chiến lược toàn cầu", "Nhóm Kinh doanh quốc tế", "Nhóm Xuất khẩu", "Nhóm Hỗ trợ quốc tế", "Nhóm Phiên dịch", "Nhóm Quản lý tài liệu", "Nhóm Quản lý tri thức",
		"Nhóm Vận hành CNTT", "Nhóm Hỗ trợ kỹ thuật", "Nhóm Chính sách bảo mật", "Nhóm Xác thực hệ thống", "Nhóm Quản lý truy cập", "Nhóm Bảo mật mạng", "Nhóm An ninh mạng", "Nhóm Bảo vệ dữ liệu", "Nhóm Quyền riêng tư", "Nhóm Đánh giá rủi ro",
		"Nhóm Trải nghiệm khách hàng", "Nhóm CX", "Nhóm Nghiên cứu người dùng", "Nhóm VOC", "Nhóm Hài lòng khách hàng", "Nhóm Quản lý cộng đồng", "Nhóm Chiến lược thương hiệu", "Nhóm Thiết kế thương hiệu", "Nhóm Khuyến mãi", "Nhóm Chiến dịch",
		"Nhóm Tiếp thị trực tuyến", "Nhóm Tiếp thị ngoại tuyến", "Nhóm Nội dung kỹ thuật số", "Nhóm Nội dung video", "Nhóm Nội dung SNS", "Nhóm Bản tin", "Nhóm Thiết kế quảng cáo", "Nhóm Minh họa", "Nhóm Thiết kế 3D", "Nhóm Thiết kế in ấn",
		"Nhóm Truyền thông nội bộ", "Nhóm Quan hệ đối ngoại", "Nhóm Quan hệ chính phủ", "Nhóm Lập kế hoạch chính sách", "Nhóm Hợp tác đại học", "Nhóm Kế hoạch nghiên cứu", "Nhóm Bằng sáng chế", "Nhóm Chứng nhận", "Nhóm Đánh giá kỹ thuật", "Nhóm Chứng nhận chất lượng",
		"Nhóm Quản lý bền vững", "Nhóm ESG", "Nhóm Môi trường và an toàn", "Nhóm An toàn công nghiệp", "Nhóm Quản lý năng lượng", "Nhóm An toàn cơ sở", "Nhóm Y tế", "Nhóm Ứng phó khẩn cấp", "Nhóm Quản lý ra vào", "Nhóm Giám sát bảo mật",
		"Nhóm Bảo mật tài liệu", "Nhóm Bảo vệ thông tin cá nhân", "Nhóm Quản lý thông tin", "Nhóm Chia sẻ tri thức", "Nhóm Học tập", "Nhóm Phát triển lãnh đạo", "Nhóm Quản lý hiệu suất", "Nhóm Phần thưởng", "Nhóm Phát triển đào tạo", "Nhóm Chiến lược nhân sự"}

	// 🔹 러시아 이름 파츠
	russianParts := []string{"Команда разработки", "Отдел планирования", "Команда дизайна", "Отдел маркетинга", "Отдел кадров", "Финансовый отдел", "Отдел продаж", "Служба поддержки клиентов", "Отдел безопасности", "Отдел контроля качества",
		"Отдел анализа данных", "Команда исследований ИИ", "Облачный отдел", "Отдел серверных операций", "Сетевой отдел", "Техническая поддержка", "Юридический отдел", "Отдел по связям с общественностью", "Отдел обучения", "Отдел стратегического планирования",
		"Команда UX", "Команда UI", "Команда мобильной разработки", "Команда веб-разработки", "Бэкенд команда", "Фронтенд команда", "Инфраструктурная команда", "QA команда", "Команда тестирования", "Продуктовая команда",
		"Отдел операций", "Отдел закупок", "Отдел управления производством", "Отдел обеспечения качества", "Отдел исследований и разработок", "Научно-исследовательский центр", "Отдел зарубежных продаж", "Отдел внутренних продаж", "Команда CS", "Отдел планирования услуг",
		"ИТ-отдел", "Технический исследовательский центр", "Отдел планирования и операций", "Контент-отдел", "Бренд-команда", "Отдел цифровой стратегии", "Рекламный отдел", "Отдел управления соцсетями", "Отдел управления партнёрами", "Отдел по работе с поставщиками",
		"Издательский отдел", "Отдел видеопроизводства", "Съёмочная команда", "Звуковая команда", "Отдел организации мероприятий", "Команда проведения событий", "Административный отдел", "Отдел управления объектами", "Отдел эксплуатации безопасности", "Команда инженеров данных",
		"Команда анализа данных", "Команда моделей ИИ", "Команда машинного обучения", "Команда робототехники", "Отдел IoT", "Команда исследований UI/UX", "Команда проектирования услуг", "Команда промышленного дизайна", "Отдел планирования контента", "CRM команда",
		"Отдел управления рекламой", "Отдел маркетинговой стратегии", "Отдел performance-маркетинга", "Исследовательский отдел", "Отдел развития талантов", "Отдел найма", "Отдел компенсаций", "Отдел трудовых отношений", "Отдел обучения и развития", "Отдел социальных льгот",
		"Бухгалтерия", "Налоговый отдел", "Бюджетный отдел", "Казначейство", "Инвестиционный отдел", "IR отдел", "Юридическая поддержка", "Отдел управления рисками", "Внутренний аудит", "Аудиторский отдел",
		"Отдел стратегических инноваций", "Отдел технологических инноваций", "Отдел нового бизнеса", "Отдел новых продуктов", "Отдел обслуживания", "Команда поддержки клиентов", "Техническая команда партнёров", "Команда мобильных сервисов", "Веб-сервис команда", "Платформенная команда",
		"Команда разработки API", "Команда DevOps", "Команда облачной инфраструктуры", "Команда разработки безопасности", "Команда шифрования", "Команда аутентификации", "Отдел эксплуатации БД", "Команда управления данными", "Команда систем ИИ", "Команда внедрения моделей",
		"Команда робототехники", "Команда проектирования систем", "Аппаратная команда", "Команда встроенного ПО", "Лаборатория обеспечения качества", "ТехЛаб", "Команда AI-решений", "Команда управления проектами", "PMO команда", "Отдел бизнес-планирования",
		"Команда новых проектов", "Команда развития партнёров", "Лаборатория сервис-дизайна", "Команда глобальной стратегии", "Отдел зарубежного бизнеса", "Отдел экспорта", "Отдел международной поддержки", "Отдел перевода", "Отдел управления документами", "Отдел управления знаниями",
		"ИТ-операционный отдел", "Служба поддержки", "Отдел политики безопасности", "Команда аутентификации систем", "Команда управления доступом", "Команда сетевой безопасности", "Отдел кибербезопасности", "Команда защиты данных", "Отдел конфиденциальности", "Команда оценки рисков",
		"Отдел клиентского опыта", "CX команда", "Команда пользовательских исследований", "Отдел VOC", "Команда удовлетворенности клиентов", "Команда управления сообществом", "Команда бренд-стратегии", "Команда бренд-дизайна", "Отдел промоушена", "Отдел кампаний",
		"Отдел онлайн-маркетинга", "Отдел офлайн-маркетинга", "Отдел цифрового контента", "Отдел видеоконтента", "Отдел контента соцсетей", "Отдел новостных рассылок", "Отдел дизайна рекламы", "Команда иллюстраторов", "Команда 3D-дизайна", "Команда печатного дизайна",
		"Отдел внутреннего PR", "Отдел внешних связей", "Отдел взаимодействия с правительством", "Отдел политики", "Отдел сотрудничества с вузами", "Отдел планирования исследований", "Отдел патентов", "Отдел сертификации", "Отдел технической оценки", "Отдел сертификации качества",
		"Отдел устойчивого развития", "ESG команда", "Отдел экологии и безопасности", "Отдел промышленной безопасности", "Отдел управления энергией", "Отдел безопасности объектов", "Медицинский отдел", "Команда экстренного реагирования", "Отдел контроля доступа", "Отдел мониторинга безопасности",
		"Отдел безопасности документов", "Отдел защиты персональных данных", "Отдел управления информацией", "Отдел обмена знаниями", "Команда обучения", "Команда развития лидерства", "Отдел управления эффективностью", "Отдел вознаграждений", "Отдел развития обучения", "Отдел стратегии HR"}

	// 🔹 일본 이름 파츠
	japaneseParts := []string{"開発チーム", "企画チーム", "デザインチーム", "人事チーム", "財務チーム", "営業チーム", "マーケティングチーム",
		"カスタマーサポートチーム", "品質管理チーム", "研究チーム", "生産チーム", "物流チーム", "購買チーム", "法務チーム", "広報チーム", "経営企画チーム", "戦略チーム", "データ分析チーム", "情報セキュリティチーム", "ネットワークチーム",
		"システム運用チーム", "クラウドチーム", "AIチーム", "機械学習チーム", "自然言語処理チーム", "画像処理チーム", "音声認識チーム", "プロダクトチーム", "サービス企画チーム", "UI/UXチーム", "フロントエンドチーム", "バックエンドチーム", "モバイルチーム",
		"テストチーム", "品質保証チーム", "監査チーム", "リスク管理チーム", "財務計画チーム", "会計チーム", "税務チーム", "資産管理チーム", "給与チーム", "教育チーム", "採用チーム", "研修チーム",
		"評価チーム", "社内文化チーム", "イベントチーム", "コミュニケーションチーム", "翻訳チーム", "通訳チーム", "コンテンツチーム", "SNS運営チーム", "ブランドチーム", "広告チーム", "市場調査チーム", "プロジェクト管理チーム", "スケジュール管理チーム", "契約管理チーム", "協力会社チーム", "サプライチェーンチーム", "製造チーム", "製品設計チーム", "安全管理チーム",
		"設備管理チーム", "環境管理チーム", "持続可能性チーム", "イノベーションチーム", "スタートアップ支援チーム", "パートナーシップチーム", "法的コンプライアンスチーム", "特許チーム",
		"著作権チーム", "データ保護チーム", "セキュリティ監査チーム", "ユーザー体験チーム", "顧客関係チーム", "ロイヤリティプログラムチーム", "会員管理チーム", "決済チーム", "請求チーム", "在庫管理チーム", "出荷チーム", "購買契約チーム", "調達チーム", "サーバーチーム", "データベースチーム", "クラウドインフラチーム", "DevOpsチーム", "API開発チーム", "内部ツールチーム", "技術支援チーム",
		"ドキュメントチーム", "ナレッジ管理チーム", "教育企画チーム", "社内研修チーム", "人材開発チーム", "人事戦略チーム", "報酬設計チーム", "人材評価チーム", "グローバルチーム", "地域マーケティングチーム",
		"海外営業チーム", "海外支社サポートチーム", "翻訳サポートチーム", "顧客データチーム", "CRMチーム", "ログ分析チーム", "AI応用チーム", "データ戦略チーム", "クラウドセキュリティチーム", "製品テストチーム", "自動化チーム", "運用最適化チーム",
		"生産管理チーム", "品質向上チーム", "コスト削減チーム", "物流最適化チーム", "供給管理チーム", "パッケージングチーム", "配送チーム", "倉庫管理チーム", "経営支援チーム", "社長室チーム", "秘書チーム", "戦略推進チーム", "企業文化チーム", "サステナビリティチーム", "環境対策チーム", "エネルギー管理チーム", "社会貢献チーム",
		"ボランティアチーム", "安全衛生チーム", "緊急対応チーム", "品質検証チーム", "監査支援チーム", "業務改善チーム", "効率化チーム", "コスト管理チーム", "プロセス改善チーム", "プロダクト戦略チーム", "市場戦略チーム", "顧客成功チーム", "技術開発チーム", "新製品開発チーム", "試作チーム", "テストエンジニアリングチーム", "サーバー運用チーム", "クラウドアーキテクチャチーム",
		"API管理チーム", "インフラモニタリングチーム", "障害対応チーム", "システム改善チーム", "UI改善チーム", "モバイルUXチーム", "ウェブUXチーム", "アプリ設計チーム", "プロトタイプチーム", "ブランドデザインチーム", "モーションデザインチーム", "イラストチーム", "写真編集チーム", "動画制作チーム", "SNSクリエイティブチーム", "デジタル広告チーム", "コピーライティングチーム", "社内広報チーム", "PRチーム", "顧客対応チーム",
		"クレーム対応チーム", "ヘルプデスクチーム", "チャットサポートチーム", "テクニカルサポートチーム", "サポート品質チーム", "教育サポートチーム", "FAQ管理チーム", "顧客満足チーム", "NPS分析チーム", "レポートチーム", "経営レポートチーム", "財務報告チーム", "業績分析チーム",
		"コスト分析チーム", "収益チーム", "予算チーム", "計画チーム", "事業開発チーム", "新規事業チーム", "提携チーム", "外部協力チーム", "研究開発チーム", "特許申請チーム", "技術文書チーム", "AI研究チーム", "製品改善チーム", "品質保証サポートチーム", "ユーザーテストチーム", "セキュリティ対策チーム", "アクセス制御チーム", "監視チーム", "脆弱性分析チーム"}

	// 🔹 중국어(간체) 이름 파츠
	chineseParts := []string{"开发团队", "企划团队", "设计团队", "市场团队", "人力资源团队", "财务团队", "销售团队", "客户支持团队", "安全团队", "质量管理团队",
		"数据分析团队", "人工智能研究团队", "云计算团队", "服务器运维团队", "网络团队", "技术支持团队", "法务团队", "公关团队", "培训团队", "战略规划团队",
		"UX团队", "UI团队", "移动开发团队", "网页开发团队", "后端团队", "前端团队", "基础设施团队", "QA团队", "测试团队", "产品团队",
		"运营团队", "采购团队", "生产管理团队", "质量保证团队", "研发团队", "研究中心", "海外销售团队", "国内销售团队", "客户服务团队", "服务策划团队",
		"IT运维团队", "技术研究中心", "策划运营团队", "内容团队", "品牌团队", "数字战略团队", "广告团队", "社交媒体运营团队", "合作伙伴管理团队", "供应商管理团队",
		"出版团队", "视频制作团队", "拍摄团队", "音效团队", "活动策划团队", "活动执行团队", "总务团队", "设施管理团队", "安全运维团队", "数据工程团队",
		"数据科学团队", "AI模型团队", "机器学习团队", "机器人团队", "物联网团队", "UI/UX研究团队", "服务设计团队", "产品设计团队", "内容策划团队", "CRM团队",
		"广告运营团队", "营销策略团队", "绩效营销团队", "研究团队", "人才发展团队", "招聘团队", "薪酬团队", "劳动关系团队", "培训运营团队", "福利团队",
		"会计团队", "税务团队", "预算团队", "资金管理团队", "投资团队", "投资者关系团队", "法务支持团队", "风险管理团队", "内部审计团队", "审计团队",
		"战略创新团队", "技术创新团队", "新业务团队", "新产品团队", "服务运营团队", "客户支持团队", "合作伙伴技术团队", "移动服务团队", "网络服务团队", "平台团队",
		"API开发团队", "DevOps团队", "云基础设施团队", "安全开发团队", "加密团队", "认证团队", "数据库运维团队", "数据管理团队", "AI系统团队", "模型部署团队",
		"机器人开发团队", "系统设计团队", "硬件团队", "固件团队", "质量保证实验室", "技术实验室", "AI解决方案团队", "项目管理团队", "PMO团队", "业务策划团队",
		"新项目团队", "合作伙伴开发团队", "服务设计实验室", "全球战略团队", "海外业务团队", "出口团队", "国际支持团队", "翻译团队", "文档管理团队", "知识管理团队",
		"IT运营团队", "帮助台团队", "安全策略团队", "认证运维团队", "访问管理团队", "网络安全团队", "网络安全团队", "数据保护团队", "隐私团队", "风险评估团队",
		"客户体验团队", "CX团队", "用户研究团队", "客户声音团队", "客户满意团队", "社区运营团队", "品牌战略团队", "品牌设计团队", "推广团队", "活动团队",
		"线上营销团队", "线下营销团队", "数字内容团队", "视频内容团队", "社交媒体内容团队", "电子报团队", "广告设计团队", "插画团队", "3D设计团队", "印刷设计团队",
		"内部公关团队", "对外合作团队", "政府合作团队", "政策策划团队", "产学合作团队", "研究策划团队", "专利团队", "认证获取团队", "技术评估团队", "质量认证团队",
		"可持续发展团队", "ESG团队", "环境与安全团队", "工业安全团队", "能源管理团队", "设施安全团队", "健康团队", "应急响应团队", "出入管理团队", "安全监控团队",
		"文档安全团队", "个人信息保护团队", "信息管理团队", "知识共享团队", "学习型组织团队", "领导力发展团队", "绩效管理团队", "薪酬奖励团队", "培训发展团队", "人力资源战略团队"}

	idx := rand.Intn(len(koreanParts))

	koreanNames := koreanParts[idx]
	englishNames := englishParts[idx]
	vietnameseNames := vietnameseParts[idx]
	russianNames := russianParts[idx]
	japaneseNames := japaneseParts[idx]
	chineseNames := chineseParts[idx]

	entity := entity.CreateMultiLangEntity{
		DeptOrg:  en.DeptOrg,
		DeptCode: en.DeptCode,
		KoLang:   koreanNames,
		EnLang:   englishNames,
		ViLang:   vietnameseNames,
		RuLang:   russianNames,
		JpLang:   japaneseNames,
		ZhLang:   chineseNames,
	}

	err := u.repository.PutWorksDeptMultiLang(ctx, entity)

	return err
}

func (u *departmentUsecase) CreateWorksDeptUser(ctx context.Context, input input.CreateWorksDeptUserInput) error {

	// 부서 조회
	depts, err := u.repository.GetDepts(ctx, input.Org)
	if err != nil {
		return err
	}

	// 매핑
	for i := 0; i < len(input.UserHashs); i++ {
		// 랜덤한 부서
		userHash := input.UserHashs[i]
		dept := depts[rand.Intn(len(depts))]
		updateHash := util.MakeUpdateHash()
		en := entity.CreateDeptUserEntity{
			UserHash:             userHash,
			DeptCode:             dept.DeptCode,
			DeptOrg:              dept.DeptOrg,
			PositionCode:         "",
			RoleCode:             "",
			IsConcurrentPosition: "N",
			UpdateHash:           updateHash,
		}
		log.Println("[CreateWorksDeptUser] user hash : ", en.UserHash, " dept : ", en.DeptCode)
		_, err := u.repository.PutDeptUser(ctx, en)
		if err != nil {
			return err
		}
	}

	return nil
}
