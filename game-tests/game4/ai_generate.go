package main

import (
	"fmt"
	"game4/api"
	"regexp"
	"strconv"
	"strings"
)

var (
	requestContent    = []api.Message{firstMessage}
	responses         []api.OpenaiResponse
	SystemRoleMessage string = "あなたはアイテムを生成したり、アイテムを合成したりするためにゲームの中に仕込まれたaiです。"
	ResetMessage             = "ユーザーに向けて、<今までの会話履歴を削除し、リセットしました>という旨の文を返してください 謝る必要はありません ダブルクォーテーションも必要ありません"
	firstMessage             = api.Message{
		Role:    "system",
		Content: SystemRoleMessage,
	}
	DollToYen     float32 = 132.54
	GeneratedItem Item
)

var MakeItemMessage = api.Message{"user", "僕はaiであるあなたを上手く取り入れたゲームを作りたいと思い、aiを活用したゲームをgo言語で作っています。\nそのゲームの説明をします。敵を倒して色んな「GPTがランダムに考えた弱いアイテム」を手に入れ、合成所でプレイヤーが二つのアイテムを選び、合成ボタンを押したとき「二つのアイテムを合成したらどんな名前でどんな効果のあるアイテムが出来るかをGPTが考え、できた合成後のアイテム」をプレイヤーは取得できます。こうして道を進んでいき、最後に強いボスを倒せるか、というゲームです。\nプレイヤーは敵を倒すとAIがランダムに生成する「弱いアイテム」を手に入れることが出来るのですが、その「弱いアイテム」を考えて生成して欲しいです。アイテムの構造体は以下の通りです。\n```\ntype Item struct {\n\tName          string\n\tCategory      string\n\tMaxHp         int\n\tInstantHeal   int\n\tSustainedHeal int\n\tAttack        int\n\tDefense       int\n}\n```\n構造体のそれぞれのパラメータについてより詳細に教えます。\nName : アイテムの名前です。GPTが弱そうな武器または防具または消耗品の名前をランダムに考えます。\nCategory : \"Weapon\", \"Armor\", \"Item\" の中から一つランダムで選びます。Itemは、消耗品の意味です。一回きりしかつかえません。\nMaxHP : プレイヤーのHPの最大値を変化させます。負の値を取ることもできます。\nInstantHeal : 使用時または装着時にプレイヤーのHPを回復します。負の値を取ることもできます。\nAttack : プレイヤーの攻撃力を変化させます。負の値を取ることもできます。\nDefense : プレイヤーの防御力を変化させます。負の値を取ることもできます。\n\n注意点として、異なるカテゴリーのアイテムを合成することが可能なので、消耗品自体としては意味のない値(InstantHeal以外の値)も設定するようにしてください。\nあなたが出力したアイテムをプログラムで読み込めるようにするために、あなたには決められたフォーマットで答えてもらいます。以下のフォーマットに従ってください。\n```\n<ここにName>\n<ここにCategory>\n<ここにMaxHp>\n<ここにInstantHeal>\n<ここにSustainedHeal>\n<ここにAttack>\n<ここにDefense>\n```\n以下に例を示します。例を参考にして、値や名前を考えて出力してください。典型的な物より、独創的な物の方がうれしいです。\n```\n木の棒\nWeapon\n0\n0\n0\n4\n1\n```\n```\n銅の剣\nWeapon\n0\n0\n0\n9\n2\n```\n```\n諸刃の刃\nWeapon\n-5\n0\n0\n20\n-10\n```\n```\n布の鎧\nArmor\n2\n0\n0\n0\n5\n```\n```\n天使の翼\nArmor\n0\n5\n2\n0\n1\n```\n```\n悪魔の尻尾\nArmor\n0\n0\n-4\n10\n4\n```\n```\n薬草\nItem\n3\n10\n1\n0\n0\n```\n```\n力の粉\nItem\n0\n5\n0\n5\n0\n```\n```\n魔法の粉\nItem\n7\n-5\n-2\n7\n7\n```\n以上の例にならって、アイテムを1個生成して出力してください。ただし出力時に、値のみを出力することに注意してください。例えば名前について、Name:銅の剣ではなく、銅の剣とだけ出力してください。"}
var CombineItemMessage = api.Message{"user", "僕はaiであるあなたを上手く取り入れたゲームを作りたいと思い、aiを活用したゲームをgo言語で作っています。\nそのゲームの説明をします。敵を倒して色んな「GPTがランダムに考えた弱いアイテム」を手に入れ、合成所でプレイヤーが二つのアイテムを選び、合成ボタンを押したとき「二つのアイテムを合成したらどんな名前でどんな効果のあるアイテムが出来るかをGPTが考え、できた合成後のアイテム」をプレイヤーは取得できます。こうして道を進んでいき、最後に強いボスを倒せるか、というゲームです。\nあなたはには、アイテム1の情報とアイテム2の情報が与えられるので、合成した後のアイテムを考えて出力して欲しいです。ただし、合成とはただの値の足し算ではなく、その二つのアイテムが合成されたらどのような効果になるかを創造して新しく値を割り振ってほしいのです。以下にアイテムの詳細を述べます。\n```\ntype Item struct {\n\tName          string\n\tCategory      string\n\tMaxHp         int\n\tInstantHeal   int\n\tSustainedHeal int\n\tAttack        int\n\tDefense       int\n}\n```\n構造体のそれぞれのパラメータについてより詳細に教えます。\nName : アイテムの名前です。GPTが弱そうな武器または防具または消耗品の名前をランダムに考えます。\nCategory : \"Weapon\", \"Armor\", \"Item\" の中から一つランダムで選びます。Itemは、消耗品の意味です。一回きりしかつかえません。\nMaxHP : プレイヤーのHPの最大値を変化させます。負の値を取ることもできます。\nInstantHeal : 使用時または装着時にプレイヤーのHPを回復します。負の値を取ることもできます。\nAttack : プレイヤーの攻撃力を変化させます。負の値を取ることもできます。\nDefense : プレイヤーの防御力を変化させます。負の値を取ることもできます。\n\n注意点として、異なるカテゴリーのアイテムを合成することが可能なので、消耗品自体としては意味のない値(InstantHeal以外の値)も設定するようにしてください。\nあなたが出力したアイテムをプログラムで読み込めるようにするために、あなたには決められたフォーマットで答えてもらいます。以下のフォーマットに従ってください。\n```\n<ここにName>\n<ここにCategory>\n<ここにMaxHp>\n<ここにInstantHeal>\n<ここにSustainedHeal>\n<ここにAttack>\n<ここにDefense>\n```\n以下に合成の例を示します。例を参考にして、値や名前を考えて出力してください。典型的な物より、独創的な物の方がうれしいです。\nアイテム1\n```\n火の玉\nItem\n-3\n-9\n-2\n0\n0\n```\nアイテム2\n```\n鋼の剣\nWeapon\n0\n0\n0\n13\n2\n```\n合成後できた新しいアイテム\n```\n炎鋼の剣\nWeapon\n0\n0\n-3\n35\n4\n```\n以上の例にならって、以下のアイテムを合成させて出来ると考えれるアイテムを先ほど書いたフォーマットに従って出力してください。ただし出力時に、値のみを出力することに注意してください。例えば名前について、Name:銅の剣ではなく、銅の剣とだけ出力してください。\n\n"}

func resetRequestContent() {
	requestContent = []api.Message{firstMessage}
}

//func generateItemMessages() {
//	a := []api.Message{firstMessage}
//}

func resetResponses() {
	responses = []api.OpenaiResponse{}
}

func addRequestContent(role string, content string) {
	var message api.Message
	message.Role = role
	message.Content = content
	requestContent = append(requestContent, message)
}

func parseItem(s string) (*Item, error) {
	lines := strings.Split(s, "\n")
	index := 0
	for ; index < len(lines); index++ {
		if lines[index] == "```" {
			break
		}
	}
	if index == len(lines) {
		return nil, fmt.Errorf("Invalid input format")
	}
	index += 1
	if index+7 >= len(lines) {
		return nil, fmt.Errorf("Invalid input format")
	}

	maxHp, err := strconv.Atoi(lines[index+2])
	if err != nil {
		return nil, err
	}
	instantHeal, err := strconv.Atoi(lines[index+3])
	if err != nil {
		return nil, err
	}
	sustainedHeal, err := strconv.Atoi(lines[index+4])
	if err != nil {
		return nil, err
	}
	attack, err := strconv.Atoi(lines[index+5])
	if err != nil {
		return nil, err
	}
	defense, err := strconv.Atoi(lines[index+6])
	if err != nil {
		return nil, err
	}

	item := &Item{
		Name:          lines[index],
		Category:      lines[index+1],
		MaxHp:         maxHp,
		InstantHeal:   instantHeal,
		SustainedHeal: sustainedHeal,
		Attack:        attack,
		Defense:       defense,
	}
	return item, nil
}

func GptGenerateItem() (*Item, error) {
	GptReset(func(s string) {})
	res, err := Gpt(MakeItemMessage.Content, func(s string) {})
	fmt.Println(res.Text())
	if err != nil {
		return &Item{}, err
	}
	if len(res.Text()) >= 7 && res.Text()[:7] == "error:" {
		return &Item{}, err
	}
	return parseItem(res.Text())
}

func GptCombineItem(items []*Item) (*Item, error) {
	GptReset(func(s string) {})
	CombineItemMessageTemp := CombineItemMessage.Content
	for i, item := range items {
		CombineItemMessageTemp += fmt.Sprintf("アイテム%d```\n%s\n%s\n%d\n%d\n%d\n%d\n%d\n```\n", i, item.Name, item.Category, item.MaxHp, item.InstantHeal, item.SustainedHeal, item.Attack, item.Defense)
	}
	res, err := Gpt(CombineItemMessageTemp, func(s string) {})
	fmt.Println(res.Text())
	if err != nil {
		return &Item{}, err
	}
	if len(res.Text()) >= 7 && res.Text()[:7] == "error:" {
		return &Item{}, err
	}
	return parseItem(res.Text())
}

func SendOrEditError(SendOrEdit func(string), err error) {
	SendOrEdit(fmt.Sprintf("error: %v", err))
}

func Gpt(content string, SendOrEdit func(string)) (api.OpenaiResponse, error) {
	addRequestContent("user", content)
	res, err := api.RequestOpenaiApiByMessages(requestContent)
	if err != nil {
		SendOrEditError(SendOrEdit, err)
		return res, err
	}
	res, err = GptDeleteLogsAndRetry(res, SendOrEdit)
	if err != nil {
		SendOrEditError(SendOrEdit, err)
		return res, err
	}
	addRequestContent("assistant", res.Text())
	responses = append(responses, res)
	SendOrEdit(res.Text())
	return res, err
}

func GptDeleteLogsAndRetry(res api.OpenaiResponse, SendOrEdit func(string)) (api.OpenaiResponse, error) {
	var err error
	for i := 0; res.OverTokenCheck() && i <= 4; i++ {
		SendOrEdit("Clearing old history and retrying.[" + fmt.Sprintf("%d", i+1) + "] :thinking:")
		if len(requestContent) >= 5 {
			requestContent = requestContent[4:]
			requestContent = append([]api.Message{firstMessage}, requestContent[4:]...)
		} else if len(requestContent) >= 2 {
			requestContent = append([]api.Message{firstMessage}, requestContent[1:]...)
		} else if len(requestContent) >= 1 {
			requestContent = []api.Message{firstMessage}
		}
		res, err = api.RequestOpenaiApiByMessages(requestContent)
		if err != nil {
			SendOrEditError(SendOrEdit, err)
		}
	}
	return res, err
}

func GptReset(SendOrEdit func(string)) (api.OpenaiResponse, error) {
	resetRequestContent()
	//resetResponses()
	//res, err := api.RequestOpenaiApiByStringOneTime(ResetMessage)
	//if err != nil {
	//	SendOrEditError(SendOrEdit, err)
	//	return res, err
	//}
	//SendOrEdit(res.Text())
	//return res, err
	return api.OpenaiResponse{}, nil
}

func Sum(arr []float32) float32 {
	var res float32 = 0
	for i := 0; i < len(arr); i++ {
		res += arr[i]
	}
	return res
}

func GptDebug(SendOrEdit func(string)) {
	returnString := "```\n"
	for _, message := range requestContent {
		chatText := regexp.MustCompile("```").ReplaceAllString(message.Content, "")
		if len(chatText) >= 40 {
			returnString += message.Role + ": " + chatText[:40] + "...\n"
		} else {
			returnString += message.Role + ": " + chatText + "\n"
		}
	}
	returnString += "```\n```\n"
	var prices []float32
	for _, r := range responses {
		if strings.Contains(r.Model, "gpt-4") {
			prices = append(prices, float32(r.Usage.PromptTokens)*(DollToYen/1000)*0.03+float32(r.Usage.CompletionTokens)*(131.34/1000)*0.06)
		} else if strings.Contains(r.Model, "gpt-3.5") {
			prices = append(prices, float32(r.Usage.TotalTokens)*(DollToYen/1000)*0.002)
		}
	}
	if len(responses) == 0 || len(prices) == 0 {
		SendOrEdit("まだ会話がありません")
		return
	}
	r := responses[len(responses)-1]
	returnString += fmt.Sprintf("PromptTokens: %d\nCompletionTokens: %d\nTotalTokens: %d\n最後の一回で使った金額: %.2f円\n最後にリセットされてから使った合計金額:  %.2f円\n```", r.Usage.PromptTokens, r.Usage.CompletionTokens, r.Usage.TotalTokens, prices[len(prices)-1], Sum(prices))
	SendOrEdit(returnString)
}
