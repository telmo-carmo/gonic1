package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"tac/gonic1/dal"

	"github.com/gin-gonic/gin"
)

func GetAllDepts(c *gin.Context) {
	myDb, _ := dal.DB.DB()

	log.Print("[INFO] GetAllDepts ")
	rows, err := myDb.Query("SELECT deptno,dname,loc,versionno,active FROM dept")
	if err != nil {
		log.Printf("[WARN] Query rows error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	myItems := []dal.Dept{}

	for rows.Next() {
		var item dal.Dept
		rows.Scan(&item.Deptno, &item.Dname, &item.Loc, &item.VersionNo, &item.Active)
		log.Printf("[INFO] Dept:%v\n", item)
		myItems = append(myItems, item)
	}
	rows.Close()

	if err = rows.Err(); err != nil {
		log.Printf("[WARN] Query rows error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": myItems})
}

func GetAllBonus(c *gin.Context) {
	myDb, _ := dal.DB.DB()

	log.Printf("[INFO] GetAllBonus ")
	rows, err := myDb.Query("SELECT ename, job, sal, comm  FROM bonus")
	if err != nil {
		log.Printf("[WARN] Query error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	myBonus := []dal.Bonus{}

	for rows.Next() {
		var tempBn dal.Bonus
		rows.Scan(&tempBn.Ename, &tempBn.Job, &tempBn.Sal, &tempBn.Comm)
		log.Printf("[INFO] Ename:%s, Job:%s, Sal:%d\n", tempBn.Ename, tempBn.Job, tempBn.Sal)
		myBonus = append(myBonus, tempBn)
	}
	rows.Close()

	if err = rows.Err(); err != nil {
		log.Printf("[WARN] Query rows error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, myBonus)
}

func GetOneBonus(c *gin.Context) {
	myDb, _ := dal.DB.DB()
	id := c.Params.ByName("id") // c.Param("id")
	uid := c.GetString("uid")
	log.Printf("[INFO] GetOneBonus id: %s (%s)", id, uid)

	var tempBn dal.Bonus
	err := myDb.QueryRow("SELECT ename, job, sal,comm FROM bonus where ename = ?", id).Scan(&tempBn.Ename, &tempBn.Job, &tempBn.Sal, &tempBn.Comm)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}
		log.Printf("[WARN] Query row error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		log.Printf("[INFO] Ename:%s, Job:%s, Sal:%d\n", tempBn.Ename, tempBn.Job, tempBn.Sal)
	}
	c.JSON(http.StatusOK, tempBn)
}

func CreateBonus(c *gin.Context) {
	// Validate input
	var bn dal.Bonus
	if err := c.ShouldBindJSON(&bn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	myDb, _ := dal.DB.DB()
	_, err := myDb.Exec("INSERT INTO bonus (ename, job, sal,comm) VALUES (?,?,?,?)", bn.Ename, bn.Job, bn.Sal, bn.Comm)
	if err != nil {
		log.Printf("[ERROR] err insert %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bn)
}

func UpdateBonus(c *gin.Context) {
	id := c.Params.ByName("id")
	// Validate input
	var bn dal.Bonus
	if err := c.ShouldBindJSON(&bn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if id != bn.Ename {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input payload and id"})
		return
	}
	myDb, _ := dal.DB.DB()
	_, err := myDb.Exec("UPDATE bonus SET job=?, sal=?,comm=?  WHERE ename = ?", bn.Job, bn.Sal, bn.Comm, bn.Ename)
	if err != nil {
		log.Printf("err %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bn)
}

func DeleteBonus(c *gin.Context) {
	id := c.Params.ByName("id")
	log.Printf("delete %s", id)
	myDb, _ := dal.DB.DB()
	var bn dal.Bonus
	err := myDb.QueryRow("SELECT ename, job, sal,comm FROM bonus where ename = ?", id).Scan(&bn.Ename, &bn.Job, &bn.Sal, &bn.Comm)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
			return
		}
		log.Printf("[WARN] Query row error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		log.Printf("[INFO] Ename:%s, Job:%s, Sal:%d\n", bn.Ename, bn.Job, bn.Sal)
	}

	_, err = myDb.Exec("DELETE FROM  bonus  WHERE ename = ?", bn.Ename)
	if err != nil {
		log.Printf("delete err %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data:": id})
}
