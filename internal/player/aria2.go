package player

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/tr1v3r/pkg/log"
)

var _ Player = (*Aria2Player)(nil)

// NewAria2Player 创建一个新的aria2播放器实例
func NewAria2Player(rpcURL, password, downloadPath string) *Aria2Player {
	return &Aria2Player{
		rpcURL:       rpcURL,
		password:     password,
		downloadPath: downloadPath,
	}
}

// Aria2Player 基于aria2的播放器实现
// 主要功能是通过aria2 RPC接口推送下载任务
// 其他控制功能无操作

type Aria2Player struct {
	rpcURL       string
	password     string
	downloadPath string
}

// Play 实现Player接口的Play方法
// 接收URL和音量，通过aria2 RPC推送下载任务
func (p *Aria2Player) Play(ctx context.Context, uri string, volume int) error {
	log.CtxDebug(ctx, "Aria2Player Play: uri=%s", uri)

	downloadOptions := make(map[string]any)
	params := make([]any, 0)

	// 对wxlivecdn.com域名的m3u8文件进行特殊处理
	if strings.Contains(uri, "wxlivecdn.com") && strings.Contains(strings.ToLower(uri), ".m3u8") {
		// 将m3u8后缀改为flv
		uri = strings.Replace(uri, ".m3u8", ".flv", 1)
		downloadOptions["header"] = `Accept: */*
Accept-Language: zh-CN,zh;q=0.9
Origin: https://channels.weixin.qq.com
Referer: https://channels.weixin.qq.com/web/pages/live?oid=zadrTig6CYU&nid=1w6yO9n_HeE&entrance_id=1009&exportkey=n_ChQIAhIQecFRCryD1CIU2OYMhtRauhL1AQIE97dBBAEAAAAAACQHKf%2FN%2BEsAAAAOpnltbLcz9gKNyK89dVj0liJOrKU22ZRuRK7096bZ89JTuwT7DUaZpiixiZ0URm3%2BfKLnXCbtazpDAJdGuYzsSWWw0sgWfJxIZN4eu%2FdKJ5OTp4PvMRzCd%2FcAhz2Iv%2BPDac2KmYXCdT5PIoUl1Z0FMAP65XE1o8QYGgCbS%2FsG2e7FMPaM%2Bq6opHJdZxq8CCVBjrXh%2FLl2aGerf%2B3Y3AVvFFRdWvq8yWI6CiHegHcBrpZIrygDdJici%2BBfHJ26Zff1sS0GWjq3VZrEhxEZbOAOw7NpwMy%2B1eG26TkbTgB3&pass_ticket=uTkjm6aVJ8%2FeY%2Ft2%2F0cugDeHru%2FoJdEVHdANSQn94GIg1AhlfmNTIriTYkX3ki98i%2FTy3mC%2Bx45OeIBLhXOBdQ%3D%3D&wx_header=0
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: cross-site
User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 NetType/WIFI MicroMessenger/7.0.20.1781(0x6700143B) WindowsWechat(0x63090a13) UnifiedPCWindowsWechat(0xf2541510) XWEB/17071 Flue`

		log.CtxDebug(ctx, "Modified uri for wxlivecdn.com: %s", uri)
	}

	// 如果有密码，添加认证信息
	if p.password != "" {
		params = append(params, "token:"+p.password)
	}

	// 从ctx中提取header信息
	// if header, ok := ctx.Value("header").(http.Header); ok {
	// 	if len(header) > 0 {
	// 		headerStr := ""
	// 		for key, values := range header {
	// 			for _, value := range values {
	// 				headerStr += key + ": " + value + "\r\n"
	// 			}
	// 		}
	// 		downloadOptions["header"] = headerStr
	// 		log.CtxDebug(ctx, "Added header to download options: %s", headerStr)
	// 	}
	// }

	// 添加下载路径（如果有）
	if p.downloadPath != "" {
		downloadOptions["dir"] = p.downloadPath
		log.CtxDebug(ctx, "Added download path: %s", p.downloadPath)
	}

	params = append(params, []string{uri}, downloadOptions)

	// 创建aria2 RPC请求
	req := aria2RPCRequest{
		JSONRPC: "2.0",
		ID:      fmt.Sprintf("%d", rand.Int()),
		Method:  "aria2.addUri",
		Params:  params,
	}

	// 发送RPC请求
	return p.sendRPCRequest(ctx, req)
}

// Pause 实现Player接口的Pause方法
// aria2播放器不支持暂停，返回nil
func (p *Aria2Player) Pause(ctx context.Context) error {
	log.CtxDebug(ctx, "Aria2Player Pause")
	return nil
}

// Stop 实现Player接口的Stop方法
// aria2播放器不支持停止，返回nil
func (p *Aria2Player) Stop(ctx context.Context) error {
	log.CtxDebug(ctx, "Aria2Player Stop")
	return nil
}

// SetVolume 实现Player接口的SetVolume方法
// aria2播放器不支持音量调节，返回nil
func (p *Aria2Player) SetVolume(ctx context.Context, v int) error {
	log.CtxDebug(ctx, "Aria2Player SetVolume: %d", v)
	return nil
}

// SetMute 实现Player接口的SetMute方法
// aria2播放器不支持静音，返回nil
func (p *Aria2Player) SetMute(ctx context.Context, m bool) error {
	log.CtxDebug(ctx, "Aria2Player SetMute: %t", m)
	return nil
}

// SetFullscreen 实现Player接口的SetFullscreen方法
// aria2播放器不支持全屏，返回nil
func (p *Aria2Player) SetFullscreen(ctx context.Context, f bool) error {
	log.CtxDebug(ctx, "Aria2Player SetFullscreen: %t", f)
	return nil
}

// SetTitle 实现Player接口的SetTitle方法
// aria2播放器不支持设置标题，返回nil
func (p *Aria2Player) SetTitle(ctx context.Context, title string) error {
	log.CtxDebug(ctx, "Aria2Player SetTitle: %s", title)
	return nil
}

// Screenshot 实现Player接口的Screenshot方法
// aria2播放器不支持截图，返回nil
func (p *Aria2Player) Screenshot(ctx context.Context, path string) error {
	log.CtxDebug(ctx, "Aria2Player Screenshot: %s", path)
	return nil
}

// SetSpeed 实现Player接口的SetSpeed方法
// aria2播放器不支持设置速度，返回nil
func (p *Aria2Player) SetSpeed(ctx context.Context, speed float64) error {
	log.CtxDebug(ctx, "Aria2Player SetSpeed: %f", speed)
	return nil
}

// GetDuration implements Player.
func (p *Aria2Player) GetDuration(ctx context.Context) (float64, error) {
	log.CtxDebug(ctx, "Aria2Player GetDuration")
	return 0, nil
}

// GetPosition implements Player.
func (p *Aria2Player) GetPosition(ctx context.Context) (float64, error) {
	log.CtxDebug(ctx, "Aria2Player GetPosition")
	return 0, nil
}

// Seek implements Player.
func (p *Aria2Player) Seek(ctx context.Context, seconds float64) error {
	log.CtxDebug(ctx, "Aria2Player Seek: seconds=%f", seconds)
	return nil
}

// aria2RPCRequest aria2 RPC请求结构
type aria2RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,omitempty"`
}

// aria2RPCResponse aria2 RPC响应结构
type aria2RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// sendRPCRequest 发送aria2 RPC请求
func (p *Aria2Player) sendRPCRequest(ctx context.Context, req aria2RPCRequest) error {
	// 序列化请求
	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal aria2 RPC request: %w", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", p.rpcURL, strings.NewReader(string(reqBody)))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// 发送HTTP请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request to aria2 RPC: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("aria2 RPC returned non-OK status: %d", resp.StatusCode)
	}

	// 解析响应
	var rpcResp aria2RPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return fmt.Errorf("failed to decode aria2 RPC response: %w", err)
	}

	// 检查RPC错误
	if rpcResp.Error != nil {
		return fmt.Errorf("aria2 RPC error: %s (code: %d)", rpcResp.Error.Message, rpcResp.Error.Code)
	}

	log.CtxDebug(ctx, "aria2 RPC request succeeded, result: %v", rpcResp.Result)
	return nil
}
