package telegram

import (
	"bytes"
	"context"
	"fmt"
	"github.com/LightningTipBot/LightningTipBot/internal"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/LightningTipBot/LightningTipBot/internal/lnbits"
	"github.com/skip2/go-qrcode"
	tb "gopkg.in/tucnak/telebot.v2"
)

func helpInvoiceUsage(ctx context.Context, errormsg string) string {
	if len(errormsg) > 0 {
		return fmt.Sprintf(Translate(ctx, "invoiceHelpText"), fmt.Sprintf("%s", errormsg))
	} else {
		return fmt.Sprintf(Translate(ctx, "invoiceHelpText"), "")
	}
}

func (bot TipBot) invoiceHandler(ctx context.Context, m *tb.Message) {
	// check and print all commands
	bot.anyTextHandler(ctx, m)
	if m.Chat.Type != tb.ChatPrivate {
		// delete message
		NewMessage(m, WithDuration(0, bot.Telegram))
		return
	}
	if len(strings.Split(m.Text, " ")) < 2 {
		bot.trySendMessage(m.Sender, helpInvoiceUsage(ctx, Translate(ctx, "invoiceEnterAmountMessage")))
		return
	}

	user := LoadUser(ctx)
	if user.Wallet == nil {
		return
	}

	userStr := GetUserStr(m.Sender)
	amount, err := decodeAmountFromCommand(m.Text)
	if err != nil {
		return
	}
	if amount > 0 {
	} else {
		bot.trySendMessage(m.Sender, helpInvoiceUsage(ctx, Translate(ctx, "invoiceValidAmountMessage")))
		return
	}

	// check for memo in command
	memo := "Powered by @LightningTipBot"
	if len(strings.Split(m.Text, " ")) > 2 {
		memo = GetMemoFromCommand(m.Text, 2)
		tag := " (@LightningTipBot)"
		memoMaxLen := 159 - len(tag)
		if len(memo) > memoMaxLen {
			memo = memo[:memoMaxLen-len(tag)]
		}
		memo = memo + tag
	}

	log.Infof("[/invoice] Creating invoice for %s of %d sat.", userStr, amount)
	// generate invoice
	invoice, err := user.Wallet.Invoice(
		lnbits.InvoiceParams{
			Out:     false,
			Amount:  int64(amount),
			Memo:    memo,
			Webhook: internal.Configuration.Lnbits.WebhookServer},
		bot.Client)
	if err != nil {
		errmsg := fmt.Sprintf("[/invoice] Could not create an invoice: %s", err)
		log.Errorln(errmsg)
		return
	}

	// create qr code
	qr, err := qrcode.Encode(invoice.PaymentRequest, qrcode.Medium, 256)
	if err != nil {
		errmsg := fmt.Sprintf("[/invoice] Failed to create QR code for invoice: %s", err)
		log.Errorln(errmsg)
		return
	}

	// send the invoice data to user
	bot.trySendMessage(m.Sender, &tb.Photo{File: tb.File{FileReader: bytes.NewReader(qr)}, Caption: fmt.Sprintf("`%s`", invoice.PaymentRequest)})
	log.Printf("[/invoice] Incvoice created. User: %s, amount: %d sat.", userStr, amount)
	return
}
