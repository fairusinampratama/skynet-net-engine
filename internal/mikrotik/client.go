package mikrotik

import (
	"fmt"
	"strconv"
	"skynet-net-engine-api/internal/models"
	
	"github.com/go-routeros/routeros"
)

type Client struct {
	Conn *routeros.Client
	Router models.Router
}

func NewClient(r models.Router) (*Client, error) {
	address := fmt.Sprintf("%s:%d", r.Host, r.Port)
	conn, err := routeros.Dial(address, r.Username, r.Password)
	if err != nil {
		return nil, err
	}
	
	return &Client{
		Conn: conn,
		Router: r,
	}, nil
}

func (c *Client) Close() {
	if c.Conn != nil {
		c.Conn.Close()
	}
}

// KeepAlive sends a lightweight command to check if connection is alive
func (c *Client) KeepAlive() error {
	_, err := c.Conn.Run("/system/identity/print")
	return err
}

func (c *Client) AddSecret(user, password, profile, localIP, remoteIP, comment string) error {
	// Construct the command
	cmd := []string{
		"/ppp/secret/add",
		"=name=" + user,
		"=password=" + password,
		"=profile=" + profile,
		"=comment=" + comment,
	}

	if localIP != "" {
		cmd = append(cmd, "=local-address="+localIP)
	}
	if remoteIP != "" {
		cmd = append(cmd, "=remote-address="+remoteIP)
	}

	_, err := c.Conn.RunArgs(cmd)
	return err
}

func (c *Client) SetSecretProfile(user, newProfile string) error {
	// Find the secret first to get ID? Or use specific query
	// RouterOS > 6.45 supports 'set [ find name=X ] profile=Y'
	
	// Doing it reliably via ID lookup is safer usually, but 'find where name=X' works
	// We will try the direct set command with a query
	
	// 1. Find ID
	res, err := c.Conn.Run("/ppp/secret/print", "?name="+user, "=.proplist=.id")
	if err != nil {
		return err
	}
	if len(res.Re) == 0 {
		return fmt.Errorf("user not found")
	}
	id := res.Re[0].Map[".id"]

	// 2. Set Profile
	_, err = c.Conn.Run("/ppp/secret/set", "=.id="+id, "=profile="+newProfile)
	return err
}

func (c *Client) AddAddressList(ip, list, comment string) error {
	_, err := c.Conn.Run(
		"/ip/firewall/address-list/add",
		"=address="+ip,
		"=list="+list,
		"=comment="+comment,
	)
	return err
}

func (c *Client) RemoveAddressList(ip, list string) error {
	// Find ID first
	res, err := c.Conn.Run("/ip/firewall/address-list/print", "?address="+ip, "?list="+list, "=.proplist=.id")
	if err != nil {
		return err
	}
	
	for _, re := range res.Re {
		id := re.Map[".id"]
		c.Conn.Run("/ip/firewall/address-list/remove", "=.id="+id)
	}
	// Ignore if not found, idempotent
	return nil
}

func (c *Client) GetActiveUsers() ([]models.ActiveUser, error) {
	res, err := c.Conn.Run("/ppp/active/print")
	if err != nil {
		return nil, err
	}

	var users []models.ActiveUser
	for _, re := range res.Re {
		users = append(users, models.ActiveUser{
			Name:     re.Map["name"],
			Address:  re.Map["address"],
			CallerID: re.Map["caller-id"],
			Uptime:   re.Map["uptime"],
			RouterID: c.Router.ID,
		})
	}
	return users, nil
}

func (c *Client) GetSystemResource() (*models.SystemResource, error) {
	res, err := c.Conn.Run("/system/resource/print")
	if err != nil {
		return nil, err
	}
	if len(res.Re) == 0 {
		return nil, fmt.Errorf("no data")
	}
	
	m := res.Re[0].Map
	// Parsing integers in production should use strconv, simplified here for MVP
	// RouterOS usually returns numbers for memory
	
	// Helper for parsing
	parseInt := func(s string) int64 {
		v, _ := strconv.ParseInt(s, 10, 64)
		return v
	}

	return &models.SystemResource{
		Uptime:      m["uptime"],
		CPU:         m["cpu-load"],
		BoardName:   m["board-name"],
		Version:     m["version"],
		TotalMemory: parseInt(m["total-memory"]),
		FreeMemory:  parseInt(m["free-memory"]),
	}, nil
}

func (c *Client) GetQueueTraffic(target string) (*models.TrafficStats, error) {
	// target can be a user name (which usually maps to a simple queue of the same name)
	// command: /queue/simple/print where name=target stats
	
	res, err := c.Conn.Run("/queue/simple/print", "?name="+target, "=.proplist=rate")
	if err != nil {
		return nil, err
	}
	if len(res.Re) == 0 {
		return nil, fmt.Errorf("queue not found")
	}
	
	// Rate comes like "rx/tx" (e.g. "1500/5000" in bits)
	rawRate := res.Re[0].Map["rate"]
	
	// Basic parsing
	var rx, tx int64
	fmt.Sscanf(rawRate, "%d/%d", &rx, &tx)
	
	return &models.TrafficStats{
		Name: target,
		RX:   rx,
		TX:   tx,
	}, nil
}

func (c *Client) RunBackup(name string) error {
	_, err := c.Conn.Run("/system/backup/save", "=name="+name)
	return err
}

