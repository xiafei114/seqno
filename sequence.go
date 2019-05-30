package seqno

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// SequenceNumber 记录表
type SequenceNumber struct {
	ID         int64  `json:"id",gorm:"auto-increment"`
	CurrentNum int64  `json:"current_Num"`
	LogicID    string `json:"logic_ID"`
	SeqFormat  string `json:"seq_Format"`
	StepNum    int64  `json:"step_Num"`
}

// SeqNo 操作代理
type SeqNo struct {
	conn           *gorm.DB
	currentNum     int64
	logicID        string
	seqFormat      string
	step           int64
	sequenceNumber *SequenceNumber
}

// NewSeqNoGenerator 初始化新代理
func NewSeqNoGenerator(db *gorm.DB, logicID string) *SeqNo {
	return &SeqNo{
		conn:       db,
		currentNum: 0,
		logicID:    logicID,
		seqFormat:  "%05d",
		step:       1,
	}
}

// InitTable 初始化表
func InitTable(db *gorm.DB) {
	db.AutoMigrate(&SequenceNumber{})
}

// Step 步长
func (s *SeqNo) Step(step int64) *SeqNo {
	s.step = step
	return s
}

// StartWith 起始数
func (s *SeqNo) StartWith(start int64) *SeqNo {
	s.currentNum = start
	return s
}

// SeqFormat 修改format
func (s *SeqNo) SeqFormat(format string) *SeqNo {
	s.seqFormat = format
	return s
}

// Next 返回序列号
func (s *SeqNo) Next() (string, error) {
	return s.next()
}

// 返回
func (s *SeqNo) next() (string, error) {
	seq := s.findCurrentSeqNumber()
	nextSeq := seq.currentNum + seq.step
	seq.currentNum = nextSeq

	seq.save()

	return seq.seqNumFormatted(), nil
}

func (s *SeqNo) save() {
	s.conn.Model(s.sequenceNumber).Update("current_Num", s.currentNum)
}

func (s *SeqNo) seqNumFormatted() string {
	return fmt.Sprintf(s.seqFormat, s.currentNum)
}

func (s *SeqNo) findCurrentSeqNumber() *SeqNo {
	var sequenceNumber SequenceNumber
	query := s.conn.First(&sequenceNumber, "logic_ID = ?", s.logicID)

	if query.Error != nil { //没有找到，新建一个
		seq := &SequenceNumber{
			CurrentNum: s.currentNum,
			LogicID:    s.logicID,
			SeqFormat:  s.seqFormat,
			StepNum:    s.step,
		}
		s.conn.Create(seq)
		s.sequenceNumber = seq
	} else {
		s.seqFormat = sequenceNumber.SeqFormat
		s.currentNum = sequenceNumber.CurrentNum
		s.step = sequenceNumber.StepNum
		s.sequenceNumber = &sequenceNumber
	}

	return s
}
