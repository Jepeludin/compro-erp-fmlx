package services

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"

	"ganttpro-backend/config"
	"ganttpro-backend/models"
)

// EmailService handles sending emails
type EmailService struct {
	config *config.Config
}

// NewEmailService creates a new EmailService
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{config: cfg}
}

// IsConfigured checks if SMTP is properly configured
func (s *EmailService) IsConfigured() bool {
	return s.config.SMTPUsername != "" && s.config.SMTPPassword != "" && s.config.SMTPFromAddr != ""
}

// GetFrontendURL returns the frontend URL for email links
func (s *EmailService) GetFrontendURL() string {
	return s.config.FrontendURL
}

// ApprovalReminderData holds data for approval reminder email templates
type ApprovalReminderData struct {
	RecipientName   string
	RecipientRole   string
	PlanID          uint
	PlanDescription string
	JobOrderNJO     string
	MachineName     string
	CreatorName     string
	SubmittedAt     string
	FrontendURL     string
	ApproveLink     string
}

// PlanApprovedData holds data for plan approval notification email
type PlanApprovedData struct {
	RecipientName   string
	PlanID          uint
	PlanNumber      string
	PlanDescription string
	JobOrderNJO     string
	MachineName     string
	ApprovedAt      string
	FrontendURL     string
}

// SendApprovalReminder sends a reminder email to an approver
func (s *EmailService) SendApprovalReminder(to *models.User, data ApprovalReminderData) error {
	if !s.IsConfigured() {
		log.Println("Email service not configured, skipping email send")
		return fmt.Errorf("email service not configured")
	}

	subject := fmt.Sprintf("[COMPRO ERP] Approval Required: Operation Plan #%d", data.PlanID)

	// Generate email body from template
	body, err := s.generateApprovalReminderHTML(data)
	if err != nil {
		return fmt.Errorf("failed to generate email body: %w", err)
	}

	return s.sendEmail(to.Email, subject, body)
}

// SendBulkApprovalReminders sends reminder emails to multiple approvers
func (s *EmailService) SendBulkApprovalReminders(approvers []*models.User, data ApprovalReminderData) []error {
	var errors []error

	for _, approver := range approvers {
		data.RecipientName = approver.Username
		data.RecipientRole = models.GetRoleDisplayName(approver.Role)

		if err := s.SendApprovalReminder(approver, data); err != nil {
			log.Printf("Failed to send email to %s: %v", approver.Email, err)
			errors = append(errors, fmt.Errorf("failed to send to %s: %w", approver.Email, err))
		} else {
			log.Printf("Email sent successfully to %s (%s)", approver.Email, approver.Role)
		}
	}

	return errors
}

// sendEmail sends an email using SMTP
func (s *EmailService) sendEmail(to, subject, body string) error {
	from := s.config.SMTPFromAddr

	// Gmail SMTP auth
	auth := smtp.PlainAuth("", s.config.SMTPUsername, s.config.SMTPPassword, s.config.SMTPHost)

	// Build email message with headers
	msg := s.buildEmailMessage(from, to, subject, body)

	// SMTP server address
	addr := fmt.Sprintf("%s:%s", s.config.SMTPHost, s.config.SMTPPort)

	// Send the email
	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// buildEmailMessage constructs the email with proper headers
func (s *EmailService) buildEmailMessage(from, to, subject, body string) []byte {
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.config.SMTPFromName, from)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var message bytes.Buffer
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	message.WriteString("\r\n")
	message.WriteString(body)

	return message.Bytes()
}

// generateApprovalReminderHTML generates HTML email content
func (s *EmailService) generateApprovalReminderHTML(data ApprovalReminderData) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            background: #ffffff;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header {
            background: linear-gradient(135deg, #1a5f7a 0%, #159895 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 24px;
        }
        .header p {
            margin: 10px 0 0;
            opacity: 0.9;
        }
        .content {
            padding: 30px;
        }
        .greeting {
            font-size: 18px;
            margin-bottom: 20px;
        }
        .info-box {
            background: #f8f9fa;
            border-left: 4px solid #159895;
            padding: 15px 20px;
            margin: 20px 0;
            border-radius: 0 4px 4px 0;
        }
        .info-row {
            display: flex;
            margin: 8px 0;
        }
        .info-label {
            font-weight: 600;
            min-width: 140px;
            color: #666;
        }
        .info-value {
            color: #333;
        }
        .cta-button {
            display: inline-block;
            background: linear-gradient(135deg, #159895 0%, #1a5f7a 100%);
            color: white !important;
            text-decoration: none;
            padding: 14px 30px;
            border-radius: 6px;
            font-weight: 600;
            margin: 20px 0;
            text-align: center;
        }
        .cta-button:hover {
            opacity: 0.9;
        }
        .footer {
            background: #f8f9fa;
            padding: 20px 30px;
            text-align: center;
            font-size: 12px;
            color: #666;
            border-top: 1px solid #eee;
        }
        .urgent-badge {
            background: #dc3545;
            color: white;
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 12px;
            font-weight: 600;
            display: inline-block;
            margin-bottom: 10px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔔 Approval Required</h1>
            <p>COMPRO ERP - Operation Plan Approval</p>
        </div>
        <div class="content">
            <p class="greeting">Hello <strong>{{.RecipientName}}</strong>,</p>
            
            <p>You have a pending operation plan that requires your approval as <strong>{{.RecipientRole}}</strong>.</p>
            
            <div class="info-box">
                <div class="info-row">
                    <span class="info-label">Plan ID:</span>
                    <span class="info-value">#{{.PlanID}}</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Job Order (NJO):</span>
                    <span class="info-value">{{.JobOrderNJO}}</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Machine:</span>
                    <span class="info-value">{{.MachineName}}</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Description:</span>
                    <span class="info-value">{{.PlanDescription}}</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Submitted By:</span>
                    <span class="info-value">{{.CreatorName}}</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Submitted At:</span>
                    <span class="info-value">{{.SubmittedAt}}</span>
                </div>
            </div>
            
            <p>Please review and approve this operation plan at your earliest convenience.</p>
            
            <center>
                <a href="{{.FrontendURL}}/operation-plans/{{.PlanID}}" class="cta-button">
                    View & Approve Plan →
                </a>
            </center>
            
            <p style="color: #666; font-size: 14px;">
                If you have any questions, please contact the creator or your supervisor.
            </p>
        </div>
        <div class="footer">
            <p>This is an automated message from ERP System.</p>
            <p>© Formulatrix - Injection Mold Engineering Division</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("approval_reminder").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *EmailService) SendPlanApprovedNotification(creator *models.User, data PlanApprovedData) error {
	if !s.IsConfigured() {
		return fmt.Errorf("email service not configured")
	}

	subject := fmt.Sprintf("[COMPRO ERP] Operation Plan #%d Approved", data.PlanID)

	body, err := s.generatePlanApprovedHTML(data)
	if err != nil {
		return fmt.Errorf("failed to generate email body: %w", err)
	}

	return s.sendEmail(creator.Email, subject, body)
}
func (s *EmailService) generatePlanApprovedHTML(data PlanApprovedData) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            line-height: 1.6;
            color: #333;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            background: #ffffff;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header {
            background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 24px;
        }
        .header p {
            margin: 10px 0 0;
            opacity: 0.9;
        }
        .content {
            padding: 30px;
        }
        .greeting {
            font-size: 18px;
            margin-bottom: 20px;
        }
        .success-badge {
            background: #d4edda;
            color: #155724;
            padding: 15px 20px;
            border-radius: 8px;
            text-align: center;
            margin: 20px 0;
            border: 1px solid #c3e6cb;
        }
        .success-badge h2 {
            margin: 0 0 5px;
            font-size: 20px;
        }
        .success-badge p {
            margin: 0;
            font-size: 14px;
        }
        .info-box {
            background: #f8f9fa;
            border-left: 4px solid #28a745;
            padding: 15px 20px;
            margin: 20px 0;
            border-radius: 0 4px 4px 0;
        }
        .info-row {
            display: flex;
            margin: 8px 0;
        }
        .info-label {
            font-weight: 600;
            min-width: 140px;
            color: #666;
        }
        .info-value {
            color: #333;
        }
        .cta-button {
            display: inline-block;
            background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
            color: white !important;
            text-decoration: none;
            padding: 14px 30px;
            border-radius: 6px;
            font-weight: 600;
            margin: 20px 0;
            text-align: center;
        }
        .next-steps {
            background: #e7f3ff;
            border-left: 4px solid #007bff;
            padding: 15px 20px;
            margin: 20px 0;
            border-radius: 0 4px 4px 0;
        }
        .next-steps h3 {
            margin: 0 0 10px;
            color: #0056b3;
        }
        .next-steps ul {
            margin: 0;
            padding-left: 20px;
        }
        .next-steps li {
            margin: 5px 0;
        }
        .footer {
            background: #f8f9fa;
            padding: 20px 30px;
            text-align: center;
            font-size: 12px;
            color: #666;
            border-top: 1px solid #eee;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Plan Approved!</h1>
            <p>COMPRO ERP - Operation Plan Notification</p>
        </div>
        <div class="content">
            <p class="greeting">Hello <strong>{{.RecipientName}}</strong>,</p>
            
            <div class="success-badge">
                <h2>🎉 Congratulations!</h2>
                <p>Your operation plan has been fully approved by all required departments.</p>
            </div>
            
            <div class="info-box">
                <div class="info-row">
                    <span class="info-label">Plan Number:</span>
                    <span class="info-value">{{.PlanNumber}}</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Plan ID:</span>
                    <span class="info-value">#{{.PlanID}}</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Job Order (NJO):</span>
                    <span class="info-value">{{.JobOrderNJO}}</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Machine:</span>
                    <span class="info-value">{{.MachineName}}</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Description:</span>
                    <span class="info-value">{{.PlanDescription}}</span>
                </div>
                <div class="info-row">
                    <span class="info-label">Approved At:</span>
                    <span class="info-value">{{.ApprovedAt}}</span>
                </div>
            </div>
            
            <div class="next-steps">
                <h3>📋 Next Steps</h3>
                <ul>
                    <li>The plan is now ready for execution</li>
                    <li>You can start the machining process using the "Start Execution" button</li>
                    <li>Remember to mark the plan as "Finished" when completed</li>
                </ul>
            </div>
            
            <center>
                <a href="{{.FrontendURL}}/operation-plans/{{.PlanID}}" class="cta-button">
                    View Plan & Start Execution →
                </a>
            </center>
        </div>
        <div class="footer">
            <p>This is an automated message from COMPRO ERP System.</p>
            <p>© Formulatrix - Injection Mold Engineering Division</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("plan_approved").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
