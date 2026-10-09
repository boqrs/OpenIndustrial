package notification

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/boqrs/nexus/email"
	"github.com/wneessen/go-mail"
)

type Service interface {
	SendAccessRequest(
		ctx context.Context,
		req AccessRequestEmail,
	) error

	SendInvitation(
		ctx context.Context,
		req InvitationEmail,
	) error
}

type AccessRequestEmail struct {
	AdminEmail string
	Name       string
	Email      string
	TenantName string
	TenantCode string
}

type InvitationEmail struct {
	Email      string
	Name       string
	TenantName string
	Token      string
}

type service struct {
	client         *email.Provider
	sender         string
	consoleBaseURL string
}

func NewService(
	client *email.Provider,
	sender string,
	consoleBaseURL string,
) Service {
	return &service{
		client:         client,
		sender:         sender,
		consoleBaseURL: strings.TrimRight(consoleBaseURL, "/"),
	}
}

func (s *service) SendAccessRequest(
	ctx context.Context,
	req AccessRequestEmail,
) error {
	if s.client == nil {
		return fmt.Errorf("email client is not configured")
	}

	if strings.TrimSpace(req.AdminEmail) == "" {
		return fmt.Errorf("admin email is empty")
	}

	msg := mail.NewMsg()

	if err := msg.From(s.sender); err != nil {
		return fmt.Errorf("set email sender: %w", err)
	}

	if err := msg.To(req.AdminEmail); err != nil {
		return fmt.Errorf("set email recipient: %w", err)
	}

	msg.Subject("[OpenIndustrial] 新用户访问申请")

	body := fmt.Sprintf(
		"您好，%s 管理员：\n\n"+
			"有一名新用户申请加入您的工厂。\n\n"+
			"姓名：%s\n"+
			"邮箱：%s\n"+
			"工厂：%s (%s)\n\n"+
			"请登录设备数字身份平台，在用户管理中审核并发送邀请。\n",
		req.TenantName,
		req.Name,
		req.Email,
		req.TenantName,
		req.TenantCode,
	)

	html := fmt.Sprintf(
		"<html><body>"+
			"<h2>新用户访问申请</h2>"+
			"<p>您好，%s 管理员：</p>"+
			"<p>有一名新用户申请加入您的工厂。</p>"+
			"<p><b>姓名：</b>%s</p>"+
			"<p><b>邮箱：</b>%s</p>"+
			"<p><b>工厂：</b>%s (%s)</p>"+
			"<p>请登录设备数字身份平台，在用户管理中审核并发送邀请。</p>"+
			"</body></html>",
		req.TenantName,
		req.Name,
		req.Email,
		req.TenantName,
		req.TenantCode,
	)

	msg.SetBodyString(mail.TypeTextPlain, body)
	msg.AddAlternativeString(mail.TypeTextHTML, html)

	if err := s.client.Get().DialAndSendWithContext(ctx, msg); err != nil {
		return fmt.Errorf("send access request email: %w", err)
	}

	return nil
}

func (s *service) SendInvitation(
	ctx context.Context,
	req InvitationEmail,
) error {
	if s.client == nil {
		return fmt.Errorf("email client is not configured")
	}

	if strings.TrimSpace(req.Email) == "" {
		return fmt.Errorf("invitation email is empty")
	}

	if strings.TrimSpace(req.Token) == "" {
		return fmt.Errorf("activation url is empty")
	}

	msg := mail.NewMsg()

	if err := msg.From(s.sender); err != nil {
		return fmt.Errorf("set email sender: %w", err)
	}

	if err := msg.To(req.Email); err != nil {
		return fmt.Errorf("set email recipient: %w", err)
	}

	msg.Subject("[OpenIndustrial] 激活您的设备数字身份平台账户")

	activationURL := BuildActivationURL(s.consoleBaseURL, req.Token)
	if activationURL == "" {
		return fmt.Errorf("console base URL is not configured")
	}

	body := fmt.Sprintf(
		"您好，%s：\n\n"+
			"您已被邀请加入 %s 的设备数字身份平台。\n\n"+
			"请打开下面的链接完成账户激活并设置密码：\n%s\n\n"+
			"该链接具有有效期，请勿转发给他人。\n",
		req.Name, req.TenantName, activationURL,
	)

	html := fmt.Sprintf(
		"<html><body>"+
			"<h2>激活您的设备数字身份平台账户</h2>"+
			"<p>您好，%s：</p>"+
			"<p>您已被邀请加入 <b>%s</b> 的设备数字身份平台。</p>"+
			"<p>请点击下面的链接完成账户激活并设置密码：</p>"+
			"<p><a href=\"%s\">激活账户</a></p>"+
			"<p>如果按钮无法打开，请复制下面的地址：</p>"+
			"<p>%s</p>"+
			"<p>该链接具有有效期，请勿转发给他人。</p>"+
			"</body></html>",
		req.Name, req.TenantName, activationURL, activationURL,
	)
	msg.SetBodyString(mail.TypeTextPlain, body)
	msg.AddAlternativeString(mail.TypeTextHTML, html)

	if err := s.client.Get().DialAndSendWithContext(ctx, msg); err != nil {
		return fmt.Errorf("send invitation email: %w", err)
	}

	return nil
}

func BuildActivationURL(baseURL string, token string) string {
    baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
    token = strings.TrimSpace(token)

    if baseURL == "" || token == "" {
        return ""
    }

    return baseURL + "/invitation/accept?token=" + url.QueryEscape(token)
}