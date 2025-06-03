package models

type ServiceUserTenant struct {
	UserHash   string `gorm:"column:user_hash;primaryKey;comment:'pk'"`
	TenantCode string `gorm:"column:tenant_code;commont:'테넌트 코드'"`
}

func (ServiceUserTenant) TableName() string {
	return "service_user_tenant"
}

// 서비스 등록 사용자 - 테넌트 코드 매핑 테이블
