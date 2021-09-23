package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"syscall"

	"github.com/howeyc/gopass"
	"github.com/polynetwork/arb-voter/config"
	"github.com/polynetwork/arb-voter/pkg/log"
	"github.com/polynetwork/arb-voter/pkg/voter"
	sdk "github.com/polynetwork/poly-go-sdk"
	"github.com/zhiqiangxu/util/signal"
)

var confFile string
var arbHeight uint64

func init() {
	flag.StringVar(&confFile, "conf", "./config.json", "configuration file path")
	flag.Uint64Var(&arbHeight, "arb", 0, "specify arb start height")
	flag.Parse()
}

func setUpPoly(polySdk *sdk.PolySdk, rpcAddr string) error {
	polySdk.NewRpcClient().SetAddress(rpcAddr)
	hdr, err := polySdk.GetHeaderByHeight(0)
	if err != nil {
		return err
	}
	polySdk.SetChainId(hdr.ChainID)
	return nil
}

func main() {
	log.InitLog(log.InfoLog, "./Log/", log.Stdout)

	conf, err := config.LoadConfig(confFile)
	if err != nil {
		log.Fatalf("LoadConfig fail:%v", err)
	}
	if arbHeight > 0 {
		conf.ForceConfig.ArbHeight = arbHeight
	}

	polySdk := sdk.NewPolySdk()
	err = setUpPoly(polySdk, conf.PolyConfig.RestURL)
	if err != nil {
		log.Fatalf("setUpPoly failed: %v", err)
	}
	wallet, err := polySdk.OpenWallet(conf.PolyConfig.WalletFile)
	if err != nil {
		log.Fatalf("polySdk.OpenWallet failed: %v", err)
	}
	fmt.Print("Enter Password: ")
	pass, err := gopass.GetPasswd()
	if err != nil {
		log.Fatalf("gopass.GetPasswd failed: %v", err)
	}
	signer, err := wallet.GetDefaultAccount(pass)
	if err != nil {
		log.Fatalf("wallet.GetDefaultAccount failed: %v", err)
	}

	log.Infof("voter %s", signer.Address.ToBase58())
	v := voter.New(polySdk, signer, conf)

	ctx, cancelFunc := context.WithCancel(context.Background())
	signal.SetupHandler(func(sig os.Signal) {
		cancelFunc()
	}, syscall.SIGINT, syscall.SIGTERM)
	v.Start(ctx)

}
