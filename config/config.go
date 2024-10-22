// // config/config.go
// package config

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// )

// type Config struct {
// 	Database struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 		Host     string `json:"host"`
// 		Port     string `json:"port"`
// 		DBName   string `json:"dbname"`
// 		Charset  string `json:"charset"`
// 	} `json:"database"`
// }

// func LoadConfig(filename string) (*Config, error) {
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	var config Config
// 	if err := json.NewDecoder(file).Decode(&config); err != nil {
// 		return nil, err
// 	}

// 	return &config, nil
// }

//	func (c *Config) GetDSN() string {
//		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
//			c.Database.Username,
//			c.Database.Password,
//			c.Database.Host,
//			c.Database.Port,
//			c.Database.DBName,
//			c.Database.Charset,
//		)
//	}
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Database struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		DBName   string `json:"dbname"`
		Charset  string `json:"charset"`
	} `json:"database"`
	JWT struct {
		SecretKey string `json:"secret_key"`
	} `json:"jwt"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	// Validate required fields
	if config.JWT.SecretKey == "" {
		return nil, fmt.Errorf("JWT secret key is required in config file")
	}

	return &config, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.Charset,
	)
}

func (c *Config) GetJWTSecret() []byte {
	return []byte(c.JWT.SecretKey)
}
