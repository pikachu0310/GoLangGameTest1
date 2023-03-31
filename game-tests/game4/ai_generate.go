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

var MakeItemMessage1 = api.Message{"user", "僕は、AIを上手く活用したゲームを作っています。\nそのゲームの説明をします。敵を倒して色んな「GPTがランダムに考えた弱いアイテム」を手に入れ、プレイヤーが二つ以上のアイテムを選び、合成ボタンを押したとき「選んだアイテムを合成したらどんな名前でどんな強さのアイテムが出来るかをGPTが考え、できた合成後のアイテム」をプレイヤーは取得できます。この二つがゲームの核です。出来たアイテムを用いて、プレイヤーは強敵を倒すことが出来るようになります。アイテムはAIによる生成だけであり、合成後のアイテムのパラメータもAIが考えるので、アイテム合成をさせた結果前より弱くなったり、またはとても強くなったりなど、どんなアイテムが出来るか分からないという楽しみがプレイヤーにあります。\n\nプレイヤーは敵を倒すとAIがランダムに生成する「弱いアイテム」を手に入れることが出来るのですが、AIであるあなたにはその「弱いアイテム」を考えて生成して欲しいです。以下にアイテムや、合成、そして出力のフォーマットに関しての詳細を述べます。\n\n## アイテムの詳細\n```\ntype Item struct {\n\tName          string\n\tCategory      string\n\tMaxHp         int\n\tInstantHeal   int\n\tSustainedHeal int\n\tAttack        int\n\tDefense       int\n}\n```\n構造体のそれぞれのパラメータについてより詳細に教えます。\n- Name : アイテムの名前です。GPTが弱そうな武器または防具または消耗品の名前をランダムに考えます。\n- Category : \"Weapon\", \"Armor\", \"Item\" の中から一つランダムで選びます。Itemは、消耗品の意味です。一回きりしかつかえません。\n- MaxHP : プレイヤーのHPの最大値を変化させます。負の値を取ることもできます。\n- InstantHeal : 使用時または装着時にプレイヤーのHPを回復します。負の値を取ることもできます。\n- Attack : プレイヤーの攻撃力を変化させます。負の値を取ることもできます。\n- Defense : プレイヤーの防御力を変化させます。負の値を取ることもできます。\n\n注意点として、異なるカテゴリーのアイテム同士であっても合成することが可能なので、消耗品自体としては意味のない値(InstantHeal以外の値)も設定するようにしてください。\n\n## 出力のフォーマット\nアイテムを出力する際は、以下のフォーマットで出力してください。\n```\n<ここにName>\n<ここにCategory>\n<ここにMaxHp>\n<ここにInstantHeal>\n<ここにSustainedHeal>\n<ここにAttack>\n<ここにDefense>\n```\nフォーマット通りの出力の例を下に書きます。\n```\n木の棒\nWeapon\n0\n0\n0\n7\n1\n```\n必ずアイテムを出力するときは上のように出力してください。\n\n## アイテム生成の詳細\nアイテム生成の注意点として、そのアイテム単体だとあまり役に立たなかったりとても弱かったりするが、合成意欲を掻き立てるようなアイテムを作ってほしいです。例えば、水や氷、草や炎といった属性っぽさがありそうなアイテム名を付けることが出来れば、プレイヤーは属性っぽさからヒントを得て面白い合成を思いつくかもしれません。例えば、氷っぽいアイテム「凍った土」と水っぽいアイテム「水鉄砲」を組み合わせれば、水を凍らせられて強いアイテムができるのではないかとか考えるかもしれません。また、属性に限らず、形容詞を付けてあげるといいかもしれません。つまり、\"弱いけど組み合わせたら強くなるかも\"なアイテム名を考えて、そのアイテムに見合う弱いパラメータを付けてください。\n\n以下にアイテムの生成例を書きます。\n## 例\n```\n石炭\nItem\n0\n-5\n0\n2\n0\n```\n```\n木の棒\nWeapon\n0\n0\n0\n4\n1\n```\n```\n銅の剣\nWeapon\n0\n0\n0\n9\n2\n```\n```\n諸刃の刃\nWeapon\n-5\n0\n0\n20\n-10\n```\n```\n布の鎧\nArmor\n2\n0\n0\n0\n5\n```\n```\n天使の翼\nArmor\n0\n5\n2\n0\n1\n```\n```\n悪魔の尻尾\nArmor\n0\n0\n-4\n10\n4\n```\n```\n薬草\nItem\n3\n10\n1\n0\n0\n```\n```\n力の粉\nItem\n0\n5\n0\n5\n0\n```\n```\n魔法の粉\nItem\n7\n-5\n-2\n7\n7\n```\n```\n黄金の十字架\nArmor\n5\n5\n0\n0\n10\n```\n```\n凍結のクロー\nWeapon\n0\n0\n0\n12\n-3\n```\n```\n火の玉\nItem\n-3\n-9\n-2\n0\n0\n```\n```\n鋼鉄のポール\nWeapon\n0\n0\n0\n14\n-4\n```\n```\n命を治す薬\nItem\n3\n20\n1\n0\n0\n```\n\n改めて、以上のようなことを踏まえて、「弱いアイテム」を1個考えて生成して、決まったフォーマットに従って出力してください。\n"}
var MakeItemMessage2 = api.Message{"assistant", "- 出力\n```\n黒い葉\nItem\n0\n-7\n2\n2\n3\n```\n\n```\n古い鎖\nArmor\n2\n0\n0\n2\n2\n```\n\n```\n魔法陣の紙\nItem\n0\n0\n5\n5\n-2\n```\n\n```\n世界樹の葉っぱ\nArmor\n4\n4\n2\n0\n3\n```\n\n```\n錆びたナイフ\nWeapon\n0\n0\n0\n11\n-2\n```"}
var MakeItemMessage3 = api.Message{"user", "完璧です！同じように、新しく「弱いアイテム」を1個考えて生成して、決まったフォーマットに従って出力してください。"}
var MakeItemMessages = []api.Message{firstMessage, MakeItemMessage1, MakeItemMessage2}

var CombineItemMessage1 = api.Message{"user", "僕は、AIを上手く活用したゲームを作っています。\nそのゲームの説明をします。敵を倒して色んな「GPTがランダムに考えた弱いアイテム」を手に入れ、プレイヤーが二つ以上のアイテムを選び、合成ボタンを押したとき「選んだアイテムを合成したらどんな名前でどんな強さのアイテムが出来るかをGPTが考え、できた合成後のアイテム」をプレイヤーは取得できます。この二つがゲームの核です。出来たアイテムを用いて、プレイヤーは強敵を倒すことが出来るようになります。アイテムはAIによる生成だけであり、合成後のアイテムのパラメータもAIが考えるので、アイテム合成をさせた結果前より弱くなったり、またはとても強くなったりなど、どんなアイテムが出来るか分からないという楽しみがプレイヤーにあります。\n\nAIであるあなたには、 \"アイテムを二つ以上受け取ってアイテムを合成させ、出来たアイテムを決まったフォーマットに従って出力する\" という事をやってほしいです。以下にアイテムや、合成、そして入力や出力のフォーマットに関しての詳細を述べます。\n\n## アイテムの詳細\n```\ntype Item struct {\n\tName          string\n\tCategory      string\n\tMaxHp         int\n\tInstantHeal   int\n\tSustainedHeal int\n\tAttack        int\n\tDefense       int\n}\n```\n構造体のそれぞれのパラメータについてより詳細に教えます。\n- Name : アイテムの名前です。GPTが弱そうな武器または防具または消耗品の名前をランダムに考えます。\n- Category : \"Weapon\", \"Armor\", \"Item\" の中から一つランダムで選びます。Itemは、消耗品の意味です。一回きりしかつかえません。\n- MaxHP : プレイヤーのHPの最大値を変化させます。負の値を取ることもできます。\n- InstantHeal : 使用時または装着時にプレイヤーのHPを回復します。負の値を取ることもできます。\n- Attack : プレイヤーの攻撃力を変化させます。負の値を取ることもできます。\n- Defense : プレイヤーの防御力を変化させます。負の値を取ることもできます。\n\n注意点として、異なるカテゴリーのアイテム同士であっても合成することが可能なので、消耗品自体としては意味のない値(InstantHeal以外の値)も設定するようにしてください。\n\n## 出力のフォーマット\nアイテムを出力する際は、以下のフォーマットで出力してください。\n```\n<ここにName>\n<ここにCategory>\n<ここにMaxHp>\n<ここにInstantHeal>\n<ここにSustainedHeal>\n<ここにAttack>\n<ここにDefense>\n```\nフォーマット通りの出力の例を下に書きます。\n```\n炎鋼の剣\nWeapon\n0\n0\n-3\n35\n4\n```\n必ずアイテムを出力するときは上のように出力してください。\n\n## 入力について\nこのプロンプトの最後に、入力を書きます。入力は以下のように書かれます。\nアイテム1\n```\n炎鋼の剣\nWeapon\n0\n0\n-3\n35\n4\n```\nアイテム2\n```\nでっかい石炭\nItem\n-5\n-10\n-1\n9\n4\n```\n\n## アイテム合成の詳細\nアイテム合成の注意点として、ただ単に値の足し算をするだけ、という事は絶対にしないでください。このアイテムとこのアイテムが組み合わさったら、いったいどういうアイテムが出来て、これくらいの強さになるのだろうか？という事を考えて作ってください。例えば、火関連のアイテムと水関連のアイテムが合成されたら打ち消しあって弱くなるかもしれません。逆に、炎っぽいアイテムと良く燃えそうなアイテムを合成したらめちゃくちゃ燃えて強くなるかもしれません。そういったアイテム同士の融合や関係性を考えて合成語のパラメータを評価してください。\n\n## 例\n- 入力\nアイテム1\n```\n炎鋼の剣\nWeapon\n0\n0\n-3\n35\n4\n```\nアイテム2\n```\nでっかい石炭\nItem\n-5\n-10\n-1\n9\n4\n```\n\n- 出力\n```\n爆炎ソード\n0\n0\n-4\n79\n5\n```\n\n改めて、以上のような事を踏まえて、以下に示す入力を受け取ってアイテムを合成させ、出来たアイテムを決まったフォーマットに従って出力してください。\n- 入力\nアイテム1\n```\n幻想の刀\nWeapon\n0\n0\n0\n14\n-4\n```\nアイテム2\n```\n火炎瓶\nItem\n0\n-10\n0\n8\n0\n```\n"}
var CombineItemMessage2 = api.Message{"assistant", "- 出力\n```\n炎幻の刀\nWeapon\n0\n-10\n0\n33\n-8\n```"}
var CombineItemMessage3 = api.Message{"user", "完璧です！同じようにして、以下の入力のアイテムの合成もお願いします。\n- 入力\n"}
var CombineItemMessages = []api.Message{firstMessage, CombineItemMessage1, CombineItemMessage2}

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
	requestContent = MakeItemMessages
	res, err := Gpt(MakeItemMessage3.Content, func(s string) {})
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
	requestContent = CombineItemMessages
	CombineItemMessageTemp := CombineItemMessage3.Content
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
