package main

import (
  "github.com/adshao/go-binance/v2"
  "log"
  "context"
  "strconv"
  "sort"
  "time"
  "fmt"
  "strings"
)

var (
  mdepLim,kinterval,kiter=1000,"1m",360
  oldwt,medwt,curwt=0.0,0.0,0.0
  basp,basl,bash,bael,baeh=0.3,0.0,0.0,10000000.0,0.0
  z,zz,d,dd=0.0,0.0,0.0
  client = binance.NewClient("","")
)

func checkFailure(e error) {if e != nil {panic(e);fmt.Println("error!")}}

func g_C() []*binance.SymbolPrice {
  p, err := client.NewListPricesService().Do(context.Background())
  checkFailure(err);return p}

type Out []string

func (list Out) has(check string) bool {
  for _,o:=range list{if o == check {return true}};return false
}

type D struct {B float64;A float64}

func mD(cLi []string) [][]D {
  mDep := make([][]D,0,len(cLi))
  for kl:=0;kl<len(cLi);kl++{
    res, err := client.NewDepthService().Symbol(cLi[kl]).Limit(mdepLim).Do(context.Background())
    checkFailure(err)
    bV:=0.0
    for h:=0;h<len(res.Bids);h++{
      bP,err:= strconv.ParseFloat(res.Bids[h].Price,64)
      checkFailure(err)
      bQ,err:=strconv.ParseFloat(res.Bids[h].Quantity,64)
      checkFailure(err)
      bV+=bP*bQ
    }
    aV:=0.0
    for r:=0;r<len(res.Asks);r++{
      aP,err:=strconv.ParseFloat(res.Asks[r].Price,64)
      checkFailure(err)
      aQ,err:= strconv.ParseFloat(res.Asks[r].Quantity,64)
      checkFailure(err)
      aV+=aP*aQ
    }
    smDep := make([]D,0,2)
    smDep = append(smDep,D{bV,aV})
    mDep = append(mDep,smDep)
    time.Sleep(1*time.Second)
  }
  return mDep
}

type DS struct {S float64;SP float64}

func sdGet(cLi []string,mdep_ [][]D) []DS {
  sss:=make([]DS,0,len(cLi))
  for _,dep:=range mdep_{
    b,a := dep[0].B,dep[0].A
    sp:=a-b
    spP:=sp/a
    valSet:=DS{sp,spP}
    sss=append(sss,valSet)
  }
  return sss
}

type ELi struct {E float64}

type _K struct {Total float64}

func get_K4Coins(_fLi []string) ([][]_K) {
  cks_ := make([][]_K,0,len(_fLi))
  for j:=0;j<len(_fLi);j++ {
    gK, err := client.NewKlinesService().
      Symbol(_fLi[j]).Interval(kinterval).Limit(kiter).Do(
        context.Background())
    if err != nil {log.Fatal(err)}
    all_K := make([]_K,0,len(gK))
    for k:=0;k<len(gK);k++{
      vM, err := strconv.ParseFloat(gK[k].Volume,64)
      checkFailure(err)
      cM, err := strconv.ParseFloat(gK[k].Close,64)
      checkFailure(err)
      valSet := _K{cM*vM,}
      all_K = append(all_K,valSet)
    }
    cks_ = append(cks_,all_K)
  }
  return cks_
}

func eG(cLi []string,allK [][]_K) []ELi{
  eee:=make([]ELi,0,len(cLi))
  for _,ack:=range allK {
    ee:=0.0
    a,b:=len(ack)*3/6,len(ack)*2/6
    for j:=0;j<len(ack);j++{
      emo:=ack[j].Total
      if j <= a-1 {
        emo *= oldwt
      } else if j >= a && j <= a+b-1 {
        emo *= medwt
      } else if j >= a + b {
        emo *= curwt
      }
      ee+=emo
    }
    eee=append(eee,ELi{ee})
  }
  return eee
}

type finalRating struct {Index int;Rate float64; Comment string}

func weightWright(
  cLi []string,
  rk []ELi,
  md []DS,
  lt string) {
    rater:=make([]finalRating,0,len(cLi))
    for r:=0;r<len(cLi);r++{
      if md[r].SP >= basp {
        if md[r].S >= basl && md[r].S < bash {
          if rk[r].E >= bael && rk[r].E < baeh {
            divergence:=rk[r].E-md[r].S
            rate:=divergence/md[r].S
            if rate < 1.0 {
              fmt.Println("Unqualified spoted")
            } else if rate >= 1.0 && rate < 2.0 {
              if md[r].S > d && rk[r].E < dd {
                rater=append(rater,finalRating{r,rate,"5xxsm5xxsl5"})
              } else if md[r].S > z && rk[r].E < zz {
                rater=append(rater,finalRating{r,rate,"4xxss6xxsm0"})
              }
            } else if rate >= 2.0 && rate < 3.0 {
              if md[r].S > z && rk[r].E < dd {
                rater=append(rater,finalRating{r,rate,"45xxsl55xss0"})
              } else if md[r].S > d && rk[r].E > dd {
                rater=append(rater,finalRating{r,rate,"6xss4xsm1"})
              }
            } else if rate >= 3.0 && rate < 4.0 {
              if md[r].S > z && rk[r].E < dd {
                rater=append(rater,finalRating{r,rate,"45xsl55ss0"})
              } else if md[r].S > d && rk[r].E > dd {
                rater=append(rater,finalRating{r,rate,"6xsm4xsl1"})
              }
            } else if rate >= 4.0 && rate < 5.0 {
              if md[r].S > z && rk[r].E < dd {
                rater=append(rater,finalRating{r,rate,"4ss6sm0"})
              } else if md[r].S > d && rk[r].E > dd {
                rater=append(rater,finalRating{r,rate,"5ss5ss5"})
              } else if md[r].S < z && rk[r].E < zz {
                rater=append(rater,finalRating{r,rate,"7sm3sl1"})
              }
            } else if rate >= 5.0 && rate < 6.0 {
              if md[r].S > z && rk[r].E < dd {
                rater=append(rater,finalRating{r,rate,"4ms5mm0"})
              } else if md[r].S > d && rk[r].E > dd {
                rater=append(rater,finalRating{r,rate,"5ms5mm5"})
              } else if md[r].S < z && rk[r].E < zz {
                rater=append(rater,finalRating{r,rate,"7ms3mm5"})
              } else if md[r].S > z && md[r].S < d {
                if rk[r].E > zz && rk[r].E < dd {
                  rater=append(rater,finalRating{r,rate,"5sl5sl1"})
                }
              }
            } else if rate >= 6.0 && rate < 7.0 {
              if md[r].S > d && rk[r].E < dd {
                rater=append(rater,finalRating{r,rate,"45ms55mm0"})
              } else if md[r].S > z && rk[r].E < zz {
                rater=append(rater,finalRating{r,rate,"5ms5mm0"})
              } else if md[r].S > d && rk[r].E > dd {
                rater=append(rater,finalRating{r,rate,"5ms5mm5"})
              } else if md[r].S < z && rk[r].E < zz {
                rater=append(rater,finalRating{r,rate,"5ms5mm0"})
              } else if md[r].S > z && md[r].S < d {
                if rk[r].E > zz && rk[r].E < dd {
                  rater=append(rater,finalRating{r,rate,"4mm6ml1"})
                }
              }
            } else if rate >= 7.0 && rate < 8.0 {
              if md[r].S > d && rk[r].E < dd {
                rater=append(rater,finalRating{r,rate,"6ml4lm0"})
              } else if md[r].S > z && rk[r].E < zz {
                rater=append(rater,finalRating{r,rate,"7ml3lm1"})
              } else if md[r].S > d && rk[r].E > dd {
                rater=append(rater,finalRating{r,rate,"5ls5lm5"})
              } else if md[r].S < z && rk[r].E < zz {
                rater=append(rater,finalRating{r,rate,"8ls2lm1"})
              } else if md[r].S > z && md[r].S < d {
                if rk[r].E > zz && rk[r].E < dd {
                  rater=append(rater,finalRating{r,rate,"55ls45ll1"})
                }
              }
            } else if rate >= 8.0 && rate < 9.0 {
              if md[r].S > d && rk[r].E < dd {
                rater=append(rater,finalRating{r,rate,"6ll4xls0"})
              } else if md[r].S > z && rk[r].E < zz {
                rater=append(rater,finalRating{r,rate,"05ll95xls0"})
              } else if md[r].S > d && rk[r].E > dd {
                rater=append(rater,finalRating{r,rate,"5xls5xlm1"})
              } else if md[r].S < z && rk[r].E < zz {
                rater=append(rater,finalRating{r,rate,"6xls4xls0"})
              } else if md[r].S > z && md[r].S < d {
                if rk[r].E > zz && rk[r].E < dd {
                  rater=append(rater,finalRating{r,rate,"8xls2xlm1"})
                }
              }
            } else if rate >= 9.0 && rate < 10.0 {
              if md[r].S < d && rk[r].E > dd {
                rater=append(rater,finalRating{r,rate,"97xxlm03xxll1"})
              } else if md[r].S > d && rk[r].E > dd {
                rater=append(rater,finalRating{r,rate,"5xxls5xxlm0"})
              } else if md[r].S > z && md[r].S < d {
                if rk[r].E > zz && rk[r].E < dd {
                  rater=append(rater,finalRating{r,rate,"65xll35xxlm0"})
                }
              } else if md[r].S < z && rk[r].E < zz {
                rater=append(rater,finalRating{r,rate,"65xll35xxlm0"})
              } else if md[r].S < z && rk[r].E > zz {
                rater=append(rater,finalRating{r,rate,"5xxlm5xxlm5"})
              }
            } else {
              if md[r].S < d && rk[r].E > dd {
                rater=append(rater,finalRating{r,rate,"35xxxxls65xxxxxxlm0"})
              } else if md[r].S > d && rk[r].E > dd {
                rater=append(rater,finalRating{r,rate,"35xxll65xxxlm0"})
              } else if md[r].S > z && md[r].S < d {
                if rk[r].E > zz && rk[r].E < dd {
                  rater=append(rater,finalRating{r,rate,"4xxlm6xxlm0"})
                }
              } else if md[r].S < z && rk[r].E < zz {
                rater=append(rater,finalRating{r,rate,"3xxlm7xxxls0"})
              } else if md[r].S < z && rk[r].E > zz {
                rater=append(rater,finalRating{r,rate,"5xxxlm5xxxxxxll0"})
              }
            }
          }
        }
      }
    }
    sort.Slice(rater,func(i,j int)bool{
      return rater[i].Rate > rater[j].Rate
    })
    if len(rater) > 0 {
      msg_:= "short term trading"
      msg_+= "\n"+lt
      for _,f:=range rater{
        msg_+="\n"+cLi[f.Index][:len(cLi[f.Index])-4]+"\n"+fmt.Sprintf("%f",f.Rate)+"\n"
        msg_+=fmt.Sprintf("%f",rk[f.Index].E)+"\n"+fmt.Sprintf("%f",md[f.Index].S)
        msg_+="\n"+fmt.Sprintf("%f",md[f.Index].SP)+"\n"+fmt.Sprintf("%s",f.Comment)+"\n"
      }
      fmt.Println(msg_)
    } else {
      fmt.Println("nothing")
    }
}

func main() {
  fmt.Println("current marked time: ",fmt.Sprintf("%s",time.Now()))
  Li:=g_C()
  LiName:=[]string{}
  for _,d := range Li {
    if strings.HasSuffix(d.Symbol,"USDT") == true {
      if strings.Contains(d.Symbol,"BULLBUSD") != true {
        if strings.Contains(d.Symbol,"BEARBUSD") != true {
          if strings.Contains(d.Symbol,"BULLUSDT") != true {
            if strings.Contains(d.Symbol,"BEARUSDT") != true {
              if strings.Contains(d.Symbol,"UPBUSD") != true {
                if strings.Contains(d.Symbol,"DOWNBUSD") != true {
                  if strings.Contains(d.Symbol,"UPUSDT") != true {
                    if strings.Contains(d.Symbol,"DOWNUSDT") != true {
                      cp,err:=strconv.ParseFloat(d.Price,64);checkFailure(err)
                      if cp <= 1000000.0 && cp >= 1.001 {
                        LiName=append(LiName,fmt.Sprintf("%s",d.Symbol))
                      }
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }
  _fLiK := get_K4Coins(LiName)
  emoS:=eG(LiName,_fLiK)
  mD_ := mD(LiName)
  sdS:=sdGet(LiName,mD_)
  weightWright(LiName,emoS,sdS,fmt.Sprintf("%s",time.Now()))
}
