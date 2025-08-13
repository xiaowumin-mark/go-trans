package youdao

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bytedance/sonic"
)

/*
	中英 en
	中法 fr
	中韩 ko
	中日 ja
*/

type YoudaoResp struct {
	Raw    string
	Parsed YoudaoResult
}

type YoudaoResult struct {
	Input string `json:"input"`
	Le    string `json:"le"`
	Meta  struct {
		Input           string   `json:"input"`
		GuessLanguage   string   `json:"guessLanguage"`
		IsHasSimpleDict string   `json:"isHasSimpleDict"` // 为"1" 则是单词或者词组 为"0"就为句子
		Le              string   `json:"le"`
		Lang            string   `json:"lang"`
		Dicts           []string `json:"dicts"`
	} `json:"meta"`

	WebTrans struct { // 网络翻译
		WebTranslation []struct {
			Key       string `json:"key"`
			KeySpeech string `json:"key-speech"`
			Trans     []struct {
				Value   string `json:"value"` // 翻译结果
				Summary struct {
					Line []string `json:"line"` // 翻译简介
				} `json:"summary"`
				Support int    `json:"support"`
				Url     string `json:"url"`
			} `json:"trans"`
			ASome string `json:"@some"` // 有这个字段则为网络释义，否则为短语
		} `json:"web-translation"`
	} `json:"web_trans"`

	Simple struct {
		Query string `json:"query"`
		Word  []struct {
			UsPhone          string `json:"usphone"`
			UkPhone          string `json:"ukphone"`
			UkSpeech         string `json:"ukspeech"`
			RrturnPhrase     string `json:"return-phrase"`
			UsSpeech         string `json:"usspeech"`
			CollegeExamVoice struct {
				SpeechWord string `json:"speechWord"`
			} `json:"collegeExamVoice"`
			Speech string `json:"speech"`
		} `json:"word"`
	} `json:"simple"`

	Phrs struct {
		Word string `json:"word"`
		Phrs []struct {
			Headword    string `json:"headword"`
			Translation string `json:"translation"`
		} `json:"phrs"`
	} `json:"phrs"`

	Syno struct {
		Word  string `json:"word"`
		Synos []struct {
			Pos  string   `json:"pos"`
			Ws   []string `json:"ws"`
			Tran string   `json:"tran"`
		} `json:"synos"`
	} `json:"syno"`

	Collins struct {
		CollinsEntries []struct {
			Hwas     string `json:"hwas"`
			Phonetic string `json:"phonetic"`
			HeadWord string `json:"head_word"`
			Star     string `json:"star"`
			Entries  struct {
				TranEntry []struct {
					PosEntry struct {
						Pos     string `json:"pos"`
						PosTips string `json:"pos_tips"`
					} `json:"pos_entry"`
					BoxExtra string `json:"box_extra"`

					ExamSents []struct {
						Sent []struct {
							ChnSent string `json:"chn_sent"`
							EngSent string `json:"eng_sent"`
						} `json:"sent"`
					} `json:"exam_sents"`
					Tran string `json:"tran"`
				} `json:"tran_entry"`
			} `json:"entries"`
			BasicEntries struct {
				BasicEntry []struct {
					Cet      string `json:"cet"`
					HeadWord string `json:"head_word"`
				} `json:"basic_entry"`
			} `json:"basic_entries"`
		} `json:"collins_entries"`
	} `json:"collins"`

	WikipediaDigest struct {
		Summarys []struct {
			Summary string `json:"summary"`
			Key     string `json:"key"`
		} `json:"summarys"`
		Sources []struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"sources"`
	} `json:"wikipedia_digest"`

	EC struct {
		WebTrans []string `json:"web_trans"` // 和网络翻译中的value保持一致
		ExamType []string `json:"exam_type"`
		Source   struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"source"`

		/*Word struct {
			UsPhone  string `json:"usphone"`
			UkPhone  string `json:"ukphone"`
			UkSpeech string `json:"ukspeech"`
			Trs      []struct {
				Pos  string `json:"pos"`
				Tran string `json:"tran"` // 网页中的简明模块
			} `json:"trs"`
			Wfs []struct {
				Wf struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"wf"`
			} `json:"wfs"`
			ReturnPhrase string `json:"return_phrase"`
			UsSpeech     string `json:"usspeech"`
		} `json:"word"`*/
		Special []struct {
			Nat   string `json:"nat"`
			Major string `json:"major"`
		} `json:"special"`
	} `json:"ec"`

	EE struct {
		Source struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"source"`
		Word struct {
			Trs []struct {
				Pos string `json:"pos"`
				Tr  []struct {
					Tran string `json:"tran"`
				} `json:"tr"`
			} `json:"trs"`
			Speech       string `json:"speech"`
			ReturnPhrase string `json:"return_phrase"`
		} `json:"word"`
	} `json:"ee"`

	BlngSentsPart struct {
		SentenceCount int `json:"sentence-count"`
		SentencePair  []struct {
			Sentence            string `json:"sentence"`
			SentenceEng         string `json:"sentence-eng"`
			SentenceTranslation string `json:"sentence-translation"`
			SpeechSize          string `json:"speech-size"`

			AlignedWords struct {
				Src struct {
					Chars []struct {
						AS     string `json:"@s"`
						AE     string `json:"@e"`
						Aid    string `json:"@id"`
						Aligns struct {
							SC []struct {
								Aid string `json:"@id"`
							} `json:"sc"`
							TC []struct {
								Aid string `json:"@id"`
							} `json:"tc"`
						} `json:"aligns"`
					} `json:"chars"`
				} `json:"src"`
				Tran struct {
					Chars []struct {
						AS     string `json:"@s"`
						AE     string `json:"@e"`
						Aid    string `json:"@id"`
						Aligns struct {
							SC []struct {
								Aid string `json:"@id"`
							} `json:"sc"`
							TC []struct {
								Aid string `json:"@id"`
							} `json:"tc"`
						} `json:"aligns"`
					} `json:"chars"`
				} `json:"tran"`
			} `json:"aligned-words"`

			Source         string `json:"source"`
			Url            string `json:"url"`
			SentenceSpeech string `json:"sentence-speech"`
		} `json:"sentence-pair"`

		More        string `json:"more"`
		TrsClassify []struct {
			Proportion string `json:"proportion"`
			TR         string `json:"tr"`
		} `json:"trs-classify"`
	} `json:"blng_sents_part"`

	Individual struct {
		Trs []struct {
			Pos  string `json:"pos"`
			Tran string `json:"tran"`
		} `json:"trs"`
		Level    string `json:"level"`
		ExamInfo struct {
			Year             int `json:"year"`
			QuestionTypeInfo []struct {
				Time int    `json:"time"`
				Type string `json:"type"`
			} `json:"questionTypeInfo"`
			ReturnPhrase  string `json:"returnPhrase"`
			PastExamSents []struct {
				En     string `json:"en"`
				Source string `json:"source"`
				Zh     string `json:"zh"`
			} `json:"pastExamSents"`
		} `json:"examInfo"`
	} `json:"individual"`

	CollinsPrimary struct {
		Words struct {
			IndexForms []string `json:"indexforms"`
			Word       string   `json:"word"`
		} `json:"words"`
		Gramcat []struct {
			Audiourl      string `json:"audiourl"`
			Pronunciation string `json:"pronunciation"`
			Senses        []struct {
				SenseNumber string `json:"sensenumber"`
				Examples    []struct {
					Sense struct {
						Lang string `json:"lang"`
						Word string `json:"word"`
					} `json:"sense"`
					Example string `json:"example"`
				} `json:"examples"`
				Definition string `json:"definition"`
				Lang       string `json:"lang"`
				Word       string `json:"word"`
			} `json:"senses"`
			Partofspeech string `json:"partofspeech"`
			Audio        string `json:"audio"`
			Forms        []struct {
				From string `json:"from"`
			} `json:"forms"`
			Variantform []struct {
				From string `json:"from"`
				Text string `json:"text"`
			} `json:"variantform"`
		} `json:"gramcat"`
	} `json:"collins_primary"`

	AuthSentsPart struct {
		SentenceCount int    `json:"sentence-count"`
		More          string `json:"more"`
		Sent          []struct {
			Score      float64 `json:"score"`
			Speech     string  `json:"speech"`
			SpeechSize string  `json:"speech_size"`
			Source     string  `json:"source"`
			Url        string  `json:"url"`
			Foreign    string  `json:"foreign"`
		} `json:"sent"`
	} `json:"auth_sents_part"`

	MediaSentsPart struct {
		SentenceCount int    `json:"sentence-count"`
		More          string `json:"more"`
		Query         string `json:"query"`
		Sent          []struct {
			AMediatype string `json:"@mediatype"`
			Snippets   struct {
				Snippet []struct {
					StreamUrl string `json:"streamUrl"`
					Duration  string `json:"duration"`
					Swf       string `json:"swf"`
					Name      string `json:"name"`
					Source    string `json:"source"`
					Win8      string `json:"win8"`
				} `json:"snippet"`
			} `json:"snippets"`
			SpeechSize int    `json:"speech_size"`
			Eng        string `json:"eng"`
		} `json:"sent"`
	} `json:"media_sents_part"`

	ExpandEc struct {
		ReturnPhrase string `json:"return_phrase"`
		Source       struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"source"`
		Word []struct {
			TransList []struct {
				Content struct {
					DetailPos string `json:"detailPos"`
					ExamType  []struct {
						En string `json:"en"`
						Zh string `json:"zh"`
					} `json:"examType"`
					Sent []struct {
						SentOrig   string `json:"sentOri"`
						SourceType string `json:"sourceType"`
						SentSpeech string `json:"sentSpeech"`
						SentTrans  string `json:"sentTrans"`
						Source     string `json:"source"`
					} `json:"sent"`
					Trans string `json:"trans"`
				} `json:"content"`
			} `json:"transList"`
			Pos string `json:"pos"`
		} `json:"word"`
	} `json:"expand_ec"`

	Etym struct {
		Etyms struct {
			Zh []struct {
				Source string `json:"source"`
				Word   string `json:"word"`
				Value  string `json:"value"`
				Url    string `json:"url"`
				Desc   string `json:"desc"`
			} `json:"zh"`
		} `json:"etyms"`
		Word string `json:"word"`
	} `json:"etym"`

	MusicSents struct {
		SentsDat []struct {
			SongName         string `json:"songName"`
			LyricTranslation string `json:"lyricTranslation"`
			Sniger           string `json:"sniger"`
			CoverImg         string `json:"coverImg"`
			SupportCount     int    `json:"supportCount"`
			Lyric            string `json:"lyric"`
			Link             string `json:"link"`
			LyricList        []struct {
				Duration         int    `json:"duration"`
				LyricTranslation string `json:"lyricTranslation"`
				Lyric            string `json:"lyric"`
				Start            int    `json:"start"`
			} `json:"lyricList"`
			Id              string `json:"id"`
			SongId          string `json:"songId"`
			DecryptedSongId string `json:"decryptedSongId"`
			PalyUrl         string `json:"playUrl"`
		} `json:"sents_dat"`
		More bool   `json:"more"`
		Word string `json:"word"`
	} `json:"music_sents"`

	Fanyi struct {
		Input string `json:"input"`
		Tran  string `json:"tran"`
		Type  string `json:"type"`
	} `json:"fanyi"`
}

type Youdao struct {
	Query string
	Lang  string
}

func New(query, lang string) *Youdao {
	return &Youdao{
		Query: query,
		Lang:  lang,
	}
}

func (y *Youdao) Translate() (*YoudaoResp, error) {
	//https://dict.youdao.com/jsonapi_s?doctype=json&jsonversion=4

	/*
		data
		q=hello&le=en&t=2&client=web&sign=3c71569a04e3231adce6ef811c67148a&keyfrom=webdict
	*/

	urlData := url.Values{
		"q":       {y.Query},
		"le":      {y.Lang},
		"client":  {"web"},
		"sign":    {"3c71569a04e3231adce6ef811c67148a"},
		"keyfrom": {"webdict"},
		"t":       {strconv.Itoa((len(y.Query + "webdict")) % 10)}, // 计算t值
	}

	req, err := http.NewRequest("POST", "https://dict.youdao.com/jsonapi_s?doctype=json&jsonversion=4", strings.NewReader(urlData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.149 Safari/537.36")
	req.Header.Set("Origin", "https://youdao.com")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result YoudaoResult
	err = sonic.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &YoudaoResp{Raw: string(body), Parsed: result}, nil
}
