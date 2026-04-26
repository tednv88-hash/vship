package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"vship/config"
)

// SendPasswordResetEmail sends a password reset email via SMTP
func SendPasswordResetEmail(cfg *config.Config, toEmail, toName, resetURL string) error {
	if cfg.SMTP.Host == "" || cfg.SMTP.FromEmail == "" {
		return fmt.Errorf("SMTP not configured")
	}

	subject := "vShip - 密碼重設"
	htmlBody := buildPasswordResetHTML(toName, resetURL, cfg.CompanyName)
	textBody := fmt.Sprintf(
		"您好 %s,\n\n您已要求重設 vShip 帳號密碼。\n\n請點擊以下連結重設密碼（連結將於 30 分鐘後失效）：\n%s\n\n如果您沒有發出此請求，請忽略此郵件。\n\n%s",
		toName, resetURL, cfg.CompanyName,
	)

	return sendSMTP(cfg, toEmail, subject, textBody, htmlBody)
}

func buildPasswordResetHTML(name, resetURL, company string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="margin:0;padding:0;font-family:'Google Sans','Segoe UI',Roboto,sans-serif;background:#f0f4f9;">
<div style="max-width:560px;margin:40px auto;background:#fff;border-radius:16px;border:1px solid #dadce0;overflow:hidden;">
  <div style="padding:40px 32px;text-align:center;">
    <div style="width:48px;height:48px;background:linear-gradient(135deg,#05ce78,#04b06a);border-radius:12px;margin:0 auto 20px;display:flex;align-items:center;justify-content:center;">
      <span style="color:#fff;font-size:24px;">&#x1f69a;</span>
    </div>
    <h1 style="font-size:24px;font-weight:400;color:#202124;margin:0 0 8px;">密碼重設</h1>
    <p style="font-size:16px;color:#5f6368;margin:0 0 24px;">vShip 跨境集運管理系統</p>
    <p style="font-size:14px;color:#3c4043;line-height:1.6;text-align:left;">
      您好 %s，<br><br>
      我們收到了您重設密碼的請求。請點擊下方按鈕重設密碼：
    </p>
    <div style="margin:32px 0;">
      <a href="%s" style="display:inline-block;padding:12px 32px;background:linear-gradient(135deg,#05ce78,#04b06a);color:#fff;text-decoration:none;border-radius:100px;font-size:14px;font-weight:500;">
        重設密碼
      </a>
    </div>
    <p style="font-size:13px;color:#5f6368;line-height:1.6;text-align:left;">
      此連結將於 <strong>30 分鐘</strong>後失效。<br>
      如果您沒有發出此請求，請忽略此郵件，您的密碼不會被更改。
    </p>
    <p style="font-size:12px;color:#9aa0a6;margin-top:24px;border-top:1px solid #e8eaed;padding-top:16px;">
      如果按鈕無法點擊，請複製以下連結到瀏覽器：<br>
      <a href="%s" style="color:#05ce78;word-break:break-all;">%s</a>
    </p>
  </div>
  <div style="background:#f8f9fa;padding:16px 32px;text-align:center;font-size:12px;color:#9aa0a6;">
    &copy; 2025 %s
  </div>
</div>
</body>
</html>`, name, resetURL, resetURL, resetURL, company)
}

func sendSMTP(cfg *config.Config, to, subject, textBody, htmlBody string) error {
	addr := net.JoinHostPort(cfg.SMTP.Host, cfg.SMTP.Port)
	fromName := cfg.SMTP.FromName
	if fromName == "" {
		fromName = "vShip"
	}

	// Build MIME message
	boundary := "vship-boundary-2025"
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", fromName, cfg.SMTP.FromEmail))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", boundary))
	msg.WriteString("\r\n")
	// Text part
	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	msg.WriteString(textBody)
	msg.WriteString("\r\n")
	// HTML part
	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
	msg.WriteString(htmlBody)
	msg.WriteString("\r\n")
	msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	// Connect
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("dial SMTP: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, cfg.SMTP.Host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	// STARTTLS
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{ServerName: cfg.SMTP.Host}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("starttls: %w", err)
		}
	}

	// Auth
	if cfg.SMTP.User != "" {
		auth := smtp.PlainAuth("", cfg.SMTP.User, cfg.SMTP.Password, cfg.SMTP.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	// Send
	if err := client.Mail(cfg.SMTP.FromEmail); err != nil {
		return fmt.Errorf("mail from: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("rcpt to: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("data: %w", err)
	}
	if _, err := w.Write([]byte(msg.String())); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}

	return client.Quit()
}
