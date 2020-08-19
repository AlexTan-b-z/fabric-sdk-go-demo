/*
 * @Author: AlexTan
 * @GIthub: https://github.com/AlexTan-b-z
 * @Date: 2020-08-04 17:44:15
 * @LastEditors: AlexTan
 * @LastEditTime: 2020-08-12 15:10:41
 */
package main

// 根据实际需求，字段自行增删,以下只供参考
type Score struct {
	ObjectType	string	`json:"docType"`
	Name	string	`json:"Name"`		// 姓名
	Gender	string	`json:"Gender"`		// 性别
	StuID	string	`json:"StuID"`		// 学生学号
	Grade	string	`json:"Grade"`		// 年级
	Result	string	`json:"Result"`		// 成绩
	Time	string	`json:"Time"`		// 插入该数据的时间

	Historys	[]HistoryItem	// 当前Score的历史记录
}

type HistoryItem struct {				// 交易hash，也是标志ID
	TxId	string
	Score	Score
}