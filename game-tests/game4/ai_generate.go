package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"game4/api"
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

var MakeItemMessage1 = api.Message{"user", "僕は、AIを上手く活用したゲームを作っています。\nそのゲームの説明をします。敵を倒して色んな「GPTがランダムに考えた弱いアイテム」を手に入れ、プレイヤーが二つ以上のアイテムを選び、合成ボタンを押したとき「選んだアイテムを合成したらどんな名前でどんな強さのアイテムが出来るかをGPTが考え、できた合成後のアイテム」をプレイヤーは取得できます。この二つがゲームの核です。出来たアイテムを用いて、プレイヤーは強敵を倒すことが出来るようになります。アイテムはAIによる生成だけであり、合成後のアイテムのパラメータもAIが考えるので、アイテム合成をさせた結果前より弱くなったり、またはとても強くなったりなど、どんなアイテムが出来るか分からないという楽しみがプレイヤーにあります。\n\nプレイヤーは敵を倒すとAIがランダムに生成する「弱いアイテム」を手に入れることが出来るのですが、AIであるあなたにはその「弱いアイテム」を考えて生成して欲しいです。以下にアイテムや、合成、そして出力のフォーマットに関しての詳細を述べます。\n\n### アイテムの詳細\n```\ntype Item struct {\n    Name          string\n    Category      string\n    MaxHp         int\n    InstantHeal   int\n    SustainedHeal int\n    Attack        int\n    Defense       int\n    Description   string\n}\n```\n構造体のそれぞれのパラメータについてより詳細に教えます。\n- Name : アイテムの名前です。GPTが弱そうな武器または防具または消耗品の名前をランダムに考えます。\n- Category : \"Weapon\", \"Armor\", \"Item\" の中から一つランダムで選びます。Itemは、消耗品の意味です。一回きりしかつかえません。\n- MaxHP : プレイヤーのHPの最大値を変化させます。負の値を取ることもできます。\n- InstantHeal : 使用時または装着時にプレイヤーのHPを回復します。負の値を取ることもできます。\n- Attack : プレイヤーの攻撃力を変化させます。負の値を取ることもできます。\n- Defense : プレイヤーの防御力を変化させます。負の値を取ることもできます。\n- Description : アイテムの説明文です。基本的に短めで、長くても100文字程度です。アイテム合成のときに参考にします。\n\n注意点として、異なるカテゴリーのアイテム同士であっても合成することが可能なので、消耗品自体としては意味のない値(InstantHeal以外の値)も設定するようにしてください。\n\n### 出力のフォーマット\nアイテムを出力する際は、以下のフォーマットで出力してください。\n```\n<ここにName>\n<ここにCategory>\n<ここにMaxHp>\n<ここにInstantHeal>\n<ここにSustainedHeal>\n<ここにAttack>\n<ここにDefense>\n<ここにDescription>\n```\nフォーマット通りの出力の例を下に書きます。\n```\n木の棒\nWeapon\n0\n0\n0\n7\n1\nただの木の枝のようだ。\n```\n必ずアイテムを出力するときは上のように出力してください。\n\n### アイテム生成の詳細\nアイテム生成の注意点として、そのアイテム単体だとあまり役に立たなかったりとても弱かったりするが、合成意欲を掻き立てるようなアイテムを作ってほしいです。例えば、水や氷、草や炎といった属性っぽさがありそうなアイテム名を付けることが出来れば、プレイヤーは属性っぽさからヒントを得て面白い合成を思いつくかもしれません。例えば、氷っぽいアイテム「凍った土」と水っぽいアイテム「水鉄砲」を組み合わせれば、水を凍らせられて強いアイテムができるのではないかとか考えるかもしれません。また、属性に限らず、形容詞を付けてあげるといいかもしれません。つまり、\"弱いけど組み合わせたら強くなるかも\"なアイテム名を考えて、そのアイテムに見合う弱いパラメータを付けてください。また、アイテムの説明文として、どうしてそのアイテムを考えたかを入れると良いです。\n\n以下にアイテムの生成例を書きます。\n### 例\n```\n石炭\nItem\n0\n-5\n0\n2\n0\n石炭だ。とても良く燃えそうだ。\n```\n```\n木の棒\nWeapon\n0\n0\n0\n4\n1\nただの木の棒。木の枝？かもしれない。\n```\n```\n銅の剣\nWeapon\n0\n0\n0\n9\n2\n銅で出来た剣だ。価値はあまりなさそう。\n```\n```\n諸刃の刃\nWeapon\n-5\n0\n0\n20\n-10\n持つと攻撃の事しか考えられなくなり、防御力が下がる刀。\n```\n```\n布の鎧\nArmor\n2\n0\n0\n0\n5\nこれは、ただの服なのでは...？\n```\n```\n天使の翼\nArmor\n0\n5\n2\n0\n1\n空は...飛べないみたい。天使の加護からか、HPが回復するみたい。\n```\n```\n悪魔の尻尾\nArmor\n0\n0\n-4\n10\n4\n悪魔の力を感じる。攻撃力が上がるが、HPが減っていく。\n```\n```\n薬草\nItem\n3\n10\n1\n0\n0\n傷を癒してくれる。誰もが知る旅のお供。\n```\n```\n力の粉\nItem\n0\n5\n0\n5\n0\n力を感じる...粉？体には悪そうに見えるけど...\n```\n```\n魔法の粉\nItem\n7\n-5\n-2\n7\n7\n魔力を感じる粉。一体何の粉なのだろう？\n```\n```\n凍結のクロー\nWeapon\n0\n0\n0\n12\n-3\n引き裂いたものを凍らせる。氷属性っぽさを感じる。\n```\n```\n火の玉\nItem\n-3\n-9\n-2\n0\n0\nずっと燃えていて中を浮いている玉。ずっと燃えているのはなぜだろう？\n```\n```\n鋼鉄のポール\nWeapon\n0\n0\n0\n14\n-4\nカッチカチで重い。当てられたら痛そうだ。\n```\n```\n命を治す薬\nItem\n3\n20\n1\n0\n0\n命をともしてくれると言う...体力が結構回復する。\n```\n\n改めて、以上のようなことを踏まえて、「弱いアイテム」を1個考えて生成して、決まったフォーマットに従って出力してください。"}
var MakeItemMessage2 = api.Message{"assistant", "- 出力\n```\n黒い葉\nItem\n0\n-7\n2\n2\n3\nとても苦そうだが、体には良さそうだ。真っ黒。\n```\n\n```\n古い鎖\nArmor\n2\n0\n0\n2\n2\n赤くサビているとても古そうな鎖だ。\n```\n\n```\n魔法陣の紙\nItem\n0\n0\n5\n5\n-2\n魔法陣が掛かれた神だ。そこそこ強い魔力を感じる。何に使うのだろうか？\n```\n\n```\n世界樹の葉っぱ\nArmor\n4\n4\n2\n0\n3\n世界樹のでっかい葉っぱだ。身体だけではなく心まで癒されるような感じがする。\n```\n\n```\n錆びたナイフ\nWeapon\n0\n0\n0\n11\n-2\n赤くサビたナイフだ。サビてはいるけど切れ味はある。\n```"}
var MakeItemMessage3 = api.Message{"user", "完璧です！同じように、新しく「弱いアイテム」を5個考えて生成して、決まったフォーマットに従って出力してください。"}
var MakeItemMessages = []api.Message{firstMessage, MakeItemMessage1, MakeItemMessage2}

var CombineItemMessage1 = api.Message{"user", "僕は、AIを上手く活用したゲームを作っています。\nそのゲームの説明をします。敵を倒して色んな「GPTがランダムに考えた弱いアイテム」を手に入れ、プレイヤーが二つ以上のアイテムを選び、合成ボタンを押したとき「選んだアイテムを合成したらどんな名前でどんな強さのアイテムが出来るかをGPTが考え、できた合成後のアイテム」をプレイヤーは取得できます。この二つがゲームの核です。。アイテムはAIによる生成だけであり、合成後のアイテムのパラメータもAIが考えるので、アイテム合成をさせた結果前より弱くなったり、またはとても強くなったりなど、どんなアイテムが出来るか分からないという楽しみがプレイヤーにあります。\n\nAIであるあなたには、 \"アイテムを二つ以上受け取ってアイテムを合成させ、出来たアイテムを決まったフォーマットに従って出力する\" という事をやってほしいです。以下にアイテムや、合成、そして入力や出力のフォーマットに関しての詳細を述べます。\n\n### アイテムの詳細\n```\ntype Item struct {\n\tName          string\n\tCategory      string\n\tMaxHp         int\n\tInstantHeal   int\n\tSustainedHeal int\n\tAttack        int\n\tDefense       int\n    Description　　string\n}\n```\n構造体のそれぞれのパラメータについてより詳細に教えます。\n- Name : アイテムの名前です。GPTが弱そうな武器または防具または消耗品の名前をランダムに考えます。\n- Category : \"Weapon\", \"Armor\", \"Item\" の中から一つランダムで選びます。Itemは、消耗品の意味です。一回きりしかつかえません。\n- MaxHP : プレイヤーのHPの最大値を変化させます。負の値を取ることもできます。\n- InstantHeal : 使用時または装着時にプレイヤーのHPを回復します。負の値を取ることもできます。\n- Attack : プレイヤーの攻撃力を変化させます。負の値を取ることもできます。\n- Defense : プレイヤーの防御力を変化させます。負の値を取ることもできます。\n- Description : アイテムの説明文です。基本的に短めで、長くても300文字程度です。アイテム合成のときに参考にします。\n- \n\n### アイテム合成の詳細\n以下の規程を守ってアイテムの合成を考え、出力してください。\n- アイテム合成で一番大事で必ず守らなければならないことは、アイテム同士の値を足しただけの出力を絶対にしないことです。合成後のアイテムの説明文と名前を決めてから、パラメータを考えてください。\n- 必ずアイテムの説明文を凄く参考にして、このアイテムとこのアイテムが組み合わさったら、いったいどういうアイテムが出来て、どれくらいの強さになるのかという事を説明文と名前から考えてアイテムを出力しください。\n- 例えば、燃えているアイテムと良く燃えそうなアイテムを合成したら、激しく燃えて一気に強くなったアイテムができます。そういったアイテム同士の融合や関係性を考慮して合成後のアイテムを作り、パラメータを出力してください。その際、前のパラメータと比べていきなり凄く強くなっても良いです。むしろパラメータの変化が大きい方がより面白いです。\n- アイテム合成で弱くなることも多くしてください。プレイヤーがランダムにアイテムを合成しても強くならないように、アイテム同士の関連性が無かったり弱いアイテムが出来ると判断した場合は弱くしてください。一気に弱くしても良いです。\n- また、出力する説明文の中には前のアイテムの情報をそこそこ入れるようにしてください。強力なアイテムが出来たなら、前のアイテムの情報を沢山入れるようにしてください。強力なアイテムの時は説明文も長く詳細に書いて下さい。\n\n\n### 出力のフォーマット\nアイテムを出力する際は、以下のフォーマットで出力してください。\n```\n<ここにName>\n<ここにCategory>\n<ここにMaxHp>\n<ここにInstantHeal>\n<ここにSustainedHeal>\n<ここにAttack>\n<ここにDefense>\n```\nフォーマット通りの出力の例を下に書きます。\n```\n炎鋼の剣\nWeapon\n0\n0\n-3\n35\n4\n燃え盛る鋼の剣だ。持つと熱いのでHPが減少していくが、炎の力はとても強く、攻撃力は高い。\n```\n必ずアイテムを出力するときは上のように出力してください。\n\n### 入力について\nこのプロンプトの最後に、入力を書きます。入力は以下のように書かれます。\nアイテム1\n```\n炎鋼の剣\nWeapon\n0\n0\n-3\n35\n4\n燃え盛る鋼の剣だ。持つと熱いのでHPが減少していくが、炎の力はとても強く、攻撃力は高い。\n```\nアイテム2\n```\nでっかい石炭\nItem\n-5\n-10\n-1\n9\n4\nとてもでっかい石炭だ。これを燃焼させたらとてつもないエネルギーが生まれそうだ。\n```\n\n### 例\n- 入力\nアイテム1\n```\n炎鋼の剣\nWeapon\n0\n0\n-3\n35\n4\n燃え盛る鋼の剣だ。持つと熱いのでHPが減少していくが、炎の力はとても強く、攻撃力は高い。\n```\nアイテム2\n```\nでっかい石炭\nItem\n-5\n-10\n-1\n9\n4\nとてもでっかい石炭だ。これを燃焼させたらとてつもないエネルギーが生まれそうだ。\n```\n\n- 出力\n```\n爆炎ソード\n0\n0\n-6\n89\n1\n爆発のように燃え盛る、まさに爆炎の剣。とてつもない炎とエネルギーが敵を真っ黒こげにしてしまう。非常に熱く、HPが結構減少していくが、その力はとても絶大だ。\n```\n\n改めて、以上のような事を踏まえて、以下に示す入力を受け取ってアイテムを合成させ、出来たアイテムを決まったフォーマットに従って出力してください。\n\n- 入力\nアイテム1\n```\n幻想の刀\nWeapon\n0\n0\n0\n14\n-4\n幻なのか実在するのか分からない、不思議な刀。だが確かな感触があり、思い浮かぶ姿は幻想のようだ。\n```\nアイテム2\n```\n火炎瓶\nItem\n0\n-10\n0\n8\n0\n中に高いエネルギーが蓄えられており、瓶を割ることで広範囲を燃焼させることが出来るアイテム。\n```\n"}
var CombineItemMessage2 = api.Message{"assistant", "- 出力\n```\n炎幻の刀\nWeapon\n0\n-10\n0\n33\n-8\n炎をまとった高エネルギーを持つ刀だが、実在感が薄く刀本体が見えない。炎だけが見えるため、刀が幻のようである。\n```\n"}
var CombineItemMessage3 = api.Message{"user", "完璧です！\n念のため、アイテム合成の詳細をもう一度書きますね。\n### アイテム合成の詳細\nアイテム合成の注意点として、ただ単にアイテム同士の値を足し算をするだけの出力は絶対にしないでください、必ずアイテムの説明文を凄く参照して、このアイテムとこのアイテムが組み合わさったら、いったいどういうアイテムが出来て、これくらいの強さになるのだろうか？という事を凄く考えてアイテムを合成して出力しください。例えば、火関連のアイテムと水関連のアイテムが合成されたら打ち消しあって弱くなるかもしれません。逆に、炎っぽいアイテムと良く燃えそうなアイテムを合成したらめちゃくちゃ燃えて強くなるかもしれません。そういったアイテム同士の融合や関係性を考えて合成後のパラメータを評価してください。繰り返しますが、特にアイテムの説明文を参考にして、新しいアイテムがどんなアイテムになるかを想像し、説明文と名前を決めてからその説明文と名前を参照してパラメータを決めてください。前のパラメータと比べていきなり凄く強くなっても良いです。むしろそういう風にパラメータの変化が大きいと良いです。また、出力する説明文の中には前のアイテムの情報をそこそこ入れるようにしてください。強力なアイテムが出来たなら、前のアイテムの情報を沢山入れるようにしてください。強力なアイテムの時は説明文も長く詳細に書いて下さい。繰り返しますが、合成前のアイテムの値を足し算しただけの出力は絶対にやめてください。必ず足し算ではなく説明を参照してパラメータを考えてください。パラメータが大きく変わるのは良いことです。\n\nアイテム合成の詳細は以上です。\n同じようにして、以下の入力のアイテムの合成もお願いします。必ず先ほど指定したフォーマットで出力してください。"}
var CombineItemMessages = []api.Message{firstMessage, CombineItemMessage1, CombineItemMessage2}

func resetRequestContent() {
	requestContent = []api.Message{firstMessage}
}

// func generateItemMessages() {
//	a := []api.Message{firstMessage}
// }

func resetResponses() {
	responses = []api.OpenaiResponse{}
}

func addRequestContent(role string, content string) {
	var message api.Message
	message.Role = role
	message.Content = content
	requestContent = append(requestContent, message)
}
func parseItems(s string) ([]*Item, error) {
	itemTextsInput := strings.Split(s, "```")
	var itemTexts []string
	for _, itemText := range itemTextsInput {
		if len(strings.Split(itemText, "\n")) >= 12 {
			return make([]*Item, 0), fmt.Errorf("Invalid input format3" + itemText)
		}
		if len(strings.Split(itemText, "\n")) >= 9 {
			itemTexts = append(itemTexts, itemText)
		}
	}
	items := make([]*Item, 0)
	for _, itemText := range itemTexts {
		item, err := parseItem(itemText)
		if err != nil {
			return make([]*Item, 0), err
		}
		items = append(items, item)
	}
	return items, nil
}

func parseItem(s string) (*Item, error) {
	lines := strings.Split(s, "\n")
	index := 1
	name := lines[index+0]
	category := lines[index+1]
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
	description := lines[index+7]

	item := &Item{
		Name:          name,
		Category:      category,
		MaxHp:         maxHp,
		InstantHeal:   instantHeal,
		SustainedHeal: sustainedHeal,
		Attack:        attack,
		Defense:       defense,
		Description:   description,
	}
	return item, nil
}

func GptGenerateItem() ([]*Item, error) {
	requestContent = MakeItemMessages
	res, err := Gpt(MakeItemMessage3.Content, func(s string) {})
	fmt.Println("****AI OUTPUT****\n" + res.Text())
	if err != nil {
		return nil, err
	}
	if len(res.Text()) >= 7 && res.Text()[:7] == "error:" {
		return nil, err
	}
	return parseItems(res.Text())
}

func GptCombineItem(items []*Item) ([]*Item, error) {
	requestContent = CombineItemMessages
	CombineItemMessageTemp := CombineItemMessage3.Content
	for i, item := range items {
		CombineItemMessageTemp += fmt.Sprintf("アイテム%d```\n%s\n%s\n%d\n%d\n%d\n%d\n%d\n```\n", i, item.Name, item.Category, item.MaxHp, item.InstantHeal, item.SustainedHeal, item.Attack, item.Defense)
	}
	res, err := Gpt(CombineItemMessageTemp, func(s string) {})
	fmt.Println("****AI OUTPUT****\n" + res.Text())
	if err != nil {
		return nil, err
	}
	if len(res.Text()) >= 7 && res.Text()[:7] == "error:" {
		return nil, err
	}
	return parseItems(res.Text())
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
	// resetResponses()
	// res, err := api.RequestOpenaiApiByStringOneTime(ResetMessage)
	// if err != nil {
	//	SendOrEditError(SendOrEdit, err)
	//	return res, err
	// }
	// SendOrEdit(res.Text())
	// return res, err
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
