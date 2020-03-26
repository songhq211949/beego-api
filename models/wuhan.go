package models

//Wuhan 武汉VO
type Wuhan struct {
	Date        string `orm:"pk" json:"date"` //以空格区分打多个tag
	SureAdd     string `json:"sureAdd"`
	SureSum     string `json:"sureSum"`
	DieAdd      string `json:"dieAdd"`
	DieSum      string `json:"dieSum"`
	RecoveryAdd string `json:"recoveryAdd"`
	RecoverySum string `json:"recoverySum"`
}
