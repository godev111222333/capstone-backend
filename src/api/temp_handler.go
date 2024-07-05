package api

import (
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func (s *Server) updateAdminReturnURL(c *gin.Context) {
	req := struct {
		NewReturnURL string `json:"new_return_url" binding:"required"`
	}{}
	if err := c.BindJSON(&req); err != nil {
		responseCustomErr(c, -1, err)
		return
	}

	s.feCfg.AdminReturnURL = req.NewReturnURL
	file, err := os.Open(s.feCfg.Path)
	if err != nil {
		responseCustomErr(c, -1, err)
		return
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()
	encoder := yaml.NewEncoder(file)
	if err := encoder.Encode(s.feCfg); err != nil {
		responseCustomErr(c, -1, err)
		return
	}

	responseSuccess(c, gin.H{"status": "update admin return url successfully"})
}
