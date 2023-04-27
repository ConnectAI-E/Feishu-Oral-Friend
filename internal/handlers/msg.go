package handlers

import (
	"context"
	"errors"
	"fmt"
	"oral-friend/internal/initialization"
	"oral-friend/pkg/openai"

	"github.com/google/uuid"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)


type CardKind string
type CardChatType string


var (
	ClearCardKind      = CardKind("clear")            // 清空上下文
	PicModeChangeKind  = CardKind("pic_mode_change")  // 切换图片创作模式
	PicResolutionKind  = CardKind("pic_resolution")   // 图片分辨率调整
	PicTextMoreKind    = CardKind("pic_text_more")    // 重新根据文本生成图片
	PicVarMoreKind     = CardKind("pic_var_more")     // 变量图片
	RoleTagsChooseKind = CardKind("role_tags_choose") // 内置角色所属标签选择
	RoleChooseKind     = CardKind("role_choose")      // 内置角色选择
	AIModeChooseKind   = CardKind("ai_mode_choose")   // AI模式选择
)
type CardMsg struct {
	Kind      CardKind
	ChatType  CardChatType
	Value     interface{}
	SessionId string
	MsgId     string
}

type MenuOption struct {
	value string
	label string
}

var (
	GroupChatType = CardChatType("group")
	UserChatType  = CardChatType("personal")
)


// withSplitLine 用于生成分割线
func withSplitLine() larkcard.MessageCardElement {
	splitLine := larkcard.NewMessageCardHr().
		Build()
	return splitLine
}


func newSendCard(header *larkcard.MessageCardHeader, elements ...larkcard.MessageCardElement) (string, error) {
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(false).
		EnableForward(true).
		UpdateMulti(true).
		Build()
	var aElementPool []larkcard.MessageCardElement
	for _, element := range elements {
		aElementPool = append(aElementPool, element)
	}
	// 卡片消息体
	cardContent, err := larkcard.NewMessageCard().
		Config(config).
		Header(header).
		Elements(
			aElementPool,
		).
		String()
	return cardContent, err
}
func newMenu(
	placeHolder string,
	value map[string]interface{},
	options ...MenuOption,
) *larkcard.
	MessageCardEmbedSelectMenuStatic {
	var aOptionPool []*larkcard.MessageCardEmbedSelectOption
	for _, option := range options {
		aOption := larkcard.NewMessageCardEmbedSelectOption().
			Value(option.value).
			Text(larkcard.NewMessageCardPlainText().
				Content(option.label).
				Build())
		aOptionPool = append(aOptionPool, aOption)

	}
	btn := larkcard.NewMessageCardEmbedSelectMenuStatic().
		MessageCardEmbedSelectMenuStatic(larkcard.NewMessageCardEmbedSelectMenuBase().
			Options(aOptionPool).
			Placeholder(larkcard.NewMessageCardPlainText().
				Content(placeHolder).
				Build()).
			Value(value).
			Build()).
		Build()
	return btn
}

// 生成带有角色标签 btn
func withRoleTagsBtn(sessionID *string, tags ...string) larkcard.
	MessageCardElement {
	var menuOptions []MenuOption

	for _, tag := range tags {
		menuOptions = append(menuOptions, MenuOption{
			label: tag,
			value: tag,
		})
	}
	cancelMenu := newMenu("选择角色分类",
		map[string]interface{}{
			"value":     "0",
			"kind":      RoleTagsChooseKind,
			"sessionId": *sessionID,
			"msgId":     *sessionID,
		},
		menuOptions...,
	)

	actions := larkcard.NewMessageCardAction().
		Actions([]larkcard.MessageCardActionElement{cancelMenu}).
		Layout(larkcard.MessageCardActionLayoutFlow.Ptr()).
		Build()
	return actions
}

// 使用 AI 模式 Btn
func withAIModeBtn(sessionID *string, aiModeStrs []string) larkcard.MessageCardElement {
	var menuOptions []MenuOption
	for _, label := range aiModeStrs {
		menuOptions = append(menuOptions, MenuOption{
			label: label,
			value: label,
		})
	}

	cancelMenu := newMenu("选择模式",
		map[string]interface{}{
			"value":     "0",
			"kind":      AIModeChooseKind,
			"sessionId": *sessionID,
			"msgId":     *sessionID,
		},
		menuOptions...,
	)

	actions := larkcard.NewMessageCardAction().
		Actions([]larkcard.MessageCardActionElement{cancelMenu}).
		Layout(larkcard.MessageCardActionLayoutFlow.Ptr()).
		Build()
	return actions
}

// withMdAndExtraBtn 用于生成带有额外按钮的消息体
func withMdAndExtraBtn(msg string, btn *larkcard.
	MessageCardEmbedButton) larkcard.MessageCardElement {
	msg, i := processMessage(msg)
	msg = processNewLine(msg)
	if i != nil {
		return nil
	}
	mainElement := larkcard.NewMessageCardDiv().
		Fields(
			[]*larkcard.MessageCardField{
				larkcard.NewMessageCardField().
					Text(larkcard.NewMessageCardLarkMd().
						Content(msg).
						Build()).
					IsShort(true).
					Build()}).
		Extra(btn).
		Build()
	return mainElement
}


// withHeader 用于生成消息头
func withHeader(title string, color string) *larkcard.
	MessageCardHeader {
	if title == "" {
		title = "🤖️机器人提醒"
	}
	header := larkcard.NewMessageCardHeader().
		Template(color).
		Title(larkcard.NewMessageCardPlainText().
			Content(title).
			Build()).
		Build()
	return header
}

// withNote 用于生成纯文本脚注
func withNote(note string) larkcard.MessageCardElement {
	noteElement := larkcard.NewMessageCardNote().
		Elements([]larkcard.MessageCardNoteElement{larkcard.NewMessageCardPlainText().
			Content(note).
			Build()}).
		Build()
	return noteElement
}

func replyCard(ctx context.Context,
	msgId *string,
	cardContent string,
) error {
	client := initialization.GetLarkClient()
	// resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
	// 	MessageId(*msgId).
	// 	Body(larkim.NewReplyMessageReqBodyBuilder().
	// 		MsgType(larkim.MsgTypeInteractive).
	// 		Uuid(uuid.New().String()).
	// 		Content(cardContent).
	// 		Build()).
	// 	Build())

		resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			Uuid(uuid.New().String()).
			Content(cardContent).
			Build()).
		Build())


	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return errors.New(resp.Msg)
	}
	return nil
}


func withRoleBtn(sessionID *string, titles ...string) larkcard.
	MessageCardElement {
	var menuOptions []MenuOption

	for _, tag := range titles {
		menuOptions = append(menuOptions, MenuOption{
			label: tag,
			value: tag,
		})
	}
	cancelMenu := newMenu("查看内置角色",
		map[string]interface{}{
			"value":     "0",
			"kind":      RoleChooseKind,
			"sessionId": *sessionID,
			"msgId":     *sessionID,
		},
		menuOptions...,
	)

	actions := larkcard.NewMessageCardAction().
		Actions([]larkcard.MessageCardActionElement{cancelMenu}).
		Layout(larkcard.MessageCardActionLayoutFlow.Ptr()).
		Build()
	return actions
}

// withMainText 用于生成纯文本消息体
func withMainText(msg string) larkcard.MessageCardElement {
	msg, i := processMessage(msg)
	msg = cleanTextBlock(msg)
	if i != nil {
		return nil
	}
	mainElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardPlainText().
				Content(msg).
				Build()).
			IsShort(false).
			Build()}).
		Build()
	return mainElement
}

func SendRoleListCard(ctx context.Context,
	sessionId *string, msgId *string, roleTag string, roleList []string) {
	newCard, _ := newSendCard(
		withHeader("🛖 角色列表"+" - "+roleTag, larkcard.TemplateIndigo),
		withRoleBtn(sessionId, roleList...),
		withNote("提醒：选择内置场景，快速进入角色扮演模式。"))
	replyCard(ctx, msgId, newCard)
}


func replyCardWithBackId(ctx context.Context,
	msgId *string,
	cardContent string,
) (*string, error) {
	client := initialization.GetLarkClient()
	resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			Uuid(uuid.New().String()).
			Content(cardContent).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return nil, errors.New(resp.Msg)
	}

	//ctx = context.WithValue(ctx, "SendMsgId", *resp.Data.MessageId)
	//SendMsgId := ctx.Value("SendMsgId")
	//pp.Println(SendMsgId)
	return resp.Data.MessageId, nil
}


func sendMsg(ctx context.Context, msg string, chatId *string) error {
	//fmt.Println("sendMsg", msg, chatId)
	msg, i := processMessage(msg)
	if i != nil {
		return i
	}
	client := initialization.GetLarkClient()
	content := larkim.NewTextMsgBuilder().
		Text(msg).
		Build()

	//fmt.Println("content", content)

	resp, err := client.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.ReceiveIdTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(*chatId).
			Content(content).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return errors.New(resp.Msg)
	}
	return nil
}

func newSendCardWithOutHeader(
	elements ...larkcard.MessageCardElement) (string, error) {
	config := larkcard.NewMessageCardConfig().
		WideScreenMode(false).
		EnableForward(true).
		UpdateMulti(true).
		Build()
	var aElementPool []larkcard.MessageCardElement
	for _, element := range elements {
		aElementPool = append(aElementPool, element)
	}
	// 卡片消息体
	cardContent, err := larkcard.NewMessageCard().
		Config(config).
		Elements(
			aElementPool,
		).
		String()
	return cardContent, err
}
// withMainMd 用于生成markdown消息体
func withMainMd(msg string) larkcard.MessageCardElement {
	msg, i := processMessage(msg)
	msg = processNewLine(msg)
	if i != nil {
		return nil
	}
	mainElement := larkcard.NewMessageCardDiv().
		Fields([]*larkcard.MessageCardField{larkcard.NewMessageCardField().
			Text(larkcard.NewMessageCardLarkMd().
				Content(msg).
				Build()).
			IsShort(true).
			Build()}).
		Build()
	return mainElement
}

func newBtn(content string, value map[string]interface{},
	typename larkcard.MessageCardButtonType) *larkcard.
	MessageCardEmbedButton {
	btn := larkcard.NewMessageCardEmbedButton().
		Type(typename).
		Value(value).
		Text(larkcard.NewMessageCardPlainText().
			Content(content).
			Build())
	return btn
}

// 清除卡片按钮
func withClearDoubleCheckBtn(sessionID *string) larkcard.MessageCardElement {
	confirmBtn := newBtn("确认清除", map[string]interface{}{
		"value":     "1",
		"kind":      ClearCardKind,
		"chatType":  UserChatType,
		"sessionId": *sessionID,
	}, larkcard.MessageCardButtonTypeDanger,
	)
	cancelBtn := newBtn("我再想想", map[string]interface{}{
		"value":     "0",
		"kind":      ClearCardKind,
		"sessionId": *sessionID,
		"chatType":  UserChatType,
	},
		larkcard.MessageCardButtonTypeDefault)

	actions := larkcard.NewMessageCardAction().
		Actions([]larkcard.MessageCardActionElement{confirmBtn, cancelBtn}).
		Layout(larkcard.MessageCardActionLayoutBisected.Ptr()).
		Build()

	return actions
}

func sendClearCacheCheckCard(ctx context.Context,
	sessionId *string, msgId *string) {
	newCard, _ := newSendCard(
		withHeader("🆑 机器人提醒", larkcard.TemplateBlue),
		withMainMd("您确定要清除对话上下文吗？"),
		withNote("请注意，这将开始一个全新的对话，您将无法利用之前话题的历史信息"),
		withClearDoubleCheckBtn(sessionId))
	replyCard(ctx, msgId, newCard)
}

func sendSystemInstructionCard(ctx context.Context,
	sessionId *string, msgId *string, content string) {
	newCard, _ := newSendCard(
		withHeader("🥷  已进入角色扮演模式", larkcard.TemplateIndigo),
		withMainText(content),
		withNote("请注意，这将开始一个全新的对话，您将无法利用之前话题的历史信息"))
	replyCard(ctx, msgId, newCard)
}

func sendOnProcessCard(ctx context.Context,
	sessionId *string, msgId *string) (*string, error) {
	newCard, _ := newSendCardWithOutHeader(
		withNote("正在思考，请稍等..."))
	id, err := replyCardWithBackId(ctx, msgId, newCard)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func sendHelpCard(ctx context.Context,
	sessionId *string, msgId *string) {
	newCard, _ := newSendCard(
		withHeader("🎒需要帮助吗？", larkcard.TemplateBlue),
		withMainMd("**我是具备打字机效果的聊天机器人！**"),
		withSplitLine(),
		withMdAndExtraBtn(
			"** 🆑 清除话题上下文**\n文本回复 *清除* 或 */clear*",
			newBtn("立刻清除", map[string]interface{}{
				"value":     "1",
				"kind":      ClearCardKind,
				"chatType":  UserChatType,
				"sessionId": *sessionId,
			}, larkcard.MessageCardButtonTypeDanger)),
		withMainMd("🛖 **内置角色列表** \n"+" 文本回复 *角色列表* 或 */roles*"),
		withMainMd("🥷 **角色扮演模式**\n文本回复*角色扮演* 或 */system*+空格+角色信息"),
		withSplitLine(),
		withMainMd("🎒 **需要更多帮助**\n文本回复 *帮助* 或 */help*"),
	)
	replyCard(ctx, msgId, newCard)
}

func replyMsg(ctx context.Context, msg string, msgId *string) error {
	msg, i := processMessage(msg)
	if i != nil {
		return i
	}
	client := initialization.GetLarkClient()
	content := larkim.NewTextMsgBuilder().
		Text(msg).
		Build()

	resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			Uuid(uuid.New().String()).
			Content(content).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return errors.New(resp.Msg)
	}
	return nil
}

func sendBalanceCard(ctx context.Context, msgId *string,
	balance openai.BalanceResponse) {
	newCard, _ := newSendCard(
		withHeader("🎰️ 余额查询", larkcard.TemplateBlue),
		withMainMd(fmt.Sprintf("总额度: %.2f$", balance.TotalGranted)),
		withMainMd(fmt.Sprintf("已用额度: %.2f$", balance.TotalUsed)),
		withMainMd(fmt.Sprintf("可用额度: %.2f$",
			balance.TotalAvailable)),
		withNote(fmt.Sprintf("有效期: %s - %s",
			balance.EffectiveAt.Format("2006-01-02 15:04:05"),
			balance.ExpiresAt.Format("2006-01-02 15:04:05"))),
	)
	replyCard(ctx, msgId, newCard)
}

func SendRoleTagsCard(ctx context.Context,
	sessionId *string, msgId *string, roleTags []string) {
	newCard, _ := newSendCard(
		withHeader("🛖 请选择角色类别", larkcard.TemplateIndigo),
		withRoleTagsBtn(sessionId, roleTags...),
		withNote("提醒：选择角色所属分类，以便我们为您推荐更多相关角色。"))
	replyCard(ctx, msgId, newCard)
}

func SendAIModeListsCard(ctx context.Context,
	sessionId *string, msgId *string, aiModeStrs []string) {
	newCard, _ := newSendCard(
		withHeader("🤖 AI模式选择", larkcard.TemplateIndigo),
		withAIModeBtn(sessionId, aiModeStrs),
		withNote("提醒：选择内置模式，让AI更好的理解您的需求。"))
	replyCard(ctx, msgId, newCard)
}

func sendNewTopicCard(ctx context.Context,
	sessionId *string, msgId *string, content string) {
	newCard, _ := newSendCard(
		withHeader("👻️ 已开启新的话题", larkcard.TemplateBlue),
		withMainText(content),
		withNote("提醒：点击对话框参与回复，可保持话题连贯"))
	replyCard(ctx, msgId, newCard)
}