package util

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"strings"
	"time"
)

const (
	defaultSize = 15
	fakeSize    = 9999
)

// Pagination 分页
type Pagination struct {
	Size      int         `form:"size" json:"size"` //size 每页显示的数量
	Page      int         `form:"page" json:"page"` //page 页码
	Total     int64       `json:"total" form:"-"`   //total 分页的查询总的数量
	List      interface{} `json:"list" form:"-"`    //data 必须是[]interface{}指针
	Order     string      `form:"order" json:"-"`   //Order 排序字段 eg1: created_at:desc  eg2:id eg3: id:desc,updated_at
	Fields    string      `form:"fields" json:"-"`  //Fields 导出excel时候选择导出的column字段名数组, 使用英文逗号分隔
	CreatedAt []time.Time `form:"created_at[]" json:"created_at"`
	UpdatedAt []time.Time `form:"updated_at[]" json:"updated_at"`
	columns   []string    //columns fields转换的字段split之后
	offset    int         //offset 分页数据偏移量
}

func checkPtrSlice(i interface{}) bool {
	if t := reflect.TypeOf(i); t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Slice {
		return true
	}
	return false
}

func parseOrderParam(r string) []clause.OrderByColumn {
	r = strings.TrimSpace(r)
	var list []clause.OrderByColumn
	if r == "" {
		return list
	}
	for _, vv := range strings.Split(r, ",") {
		ps := strings.Split(vv, ":")
		if len(ps) == 2 {
			isDesc := strings.TrimSpace(ps[1]) == "desc"
			list = append(list, clause.OrderByColumn{Column: clause.Column{Name: strings.TrimSpace(ps[0])}, Desc: isDesc})
		}
		if len(ps) == 1 {
			list = append(list, clause.OrderByColumn{Column: clause.Column{Name: strings.TrimSpace(ps[0])}, Desc: false})
		}
	}
	return list
}

//Fetch 支持构建原始sql  可以保持软删除的一致性
//tx gorm.DB 构造的自定义查询条件
//slicePtr slice的指针引用
//如果时间到字段的计算 请自定义Gorm AfterFind 钩子函数
func (p *Pagination) Fetch(tx *gorm.DB, slicePtr interface{}) (err error) {
	if !checkPtrSlice(slicePtr) {
		return errors.New("list parameter must be a slice pointer")
	}
	p.paramClean()
	if p.Size == fakeSize {
		err = tx.Find(slicePtr).Error
		if err != nil {
			return err
		}
	} else {
		var total int64
		err = tx.Count(&total).Error
		if err != nil {
			return err
		}
		p.Total = total

		for _, ss := range parseOrderParam(p.Order) {
			tx = tx.Order(ss)
		}
		err = tx.Limit(p.Size).Offset(p.offset).Find(slicePtr).Error
		if err != nil {
			return err
		}
	}
	p.List = slicePtr
	return
}

//paramClean 分页参数清洗
func (p *Pagination) paramClean() *Pagination {
	if p.Size < 1 {
		p.Size = defaultSize
	}
	if p.Page < 1 {
		p.Page = 1
	}
	var cols []string
	for _, ssss := range strings.Split(p.Fields, ",") {
		if vs := strings.TrimSpace(ssss); vs != "" {
			cols = append(cols, vs)
		}
	}
	p.columns = cols
	p.offset = p.Size * (p.Page - 1)
	return p
}
