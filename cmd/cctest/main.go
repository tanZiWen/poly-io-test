/*
* Copyright (C) 2020 The poly network Authors
* This file is part of The poly network library.
*
* The poly network is free software: you can redistribute it and/or modify
* it under the terms of the GNU Lesser General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* The poly network is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
* GNU Lesser General Public License for more details.
* You should have received a copy of the GNU Lesser General Public License
* along with The poly network . If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"flag"
	"fmt"
	"github.com/polynetwork/cross_chain_test/chains/btc"
	"github.com/polynetwork/cross_chain_test/chains/cosmos"
	"github.com/polynetwork/cross_chain_test/chains/eth"
	"github.com/polynetwork/cross_chain_test/chains/ont"
	"github.com/polynetwork/cross_chain_test/config"
	"github.com/polynetwork/cross_chain_test/log"
	_ "github.com/polynetwork/cross_chain_test/testcase"
	"github.com/polynetwork/cross_chain_test/testframework"
	"github.com/polynetwork/poly-go-sdk"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	TestConfig string //Test config file
	TestCases  string //TestCase list in cmdline
	LoopNumber int
)

func init() {
	flag.StringVar(&TestConfig, "cfg", "./config.json", "Config of cross_chain_test")
	flag.StringVar(&TestCases, "t", "", "Test case to run. use ',' to split test case")
	flag.IntVar(&LoopNumber, "loop", 1, " the number the whole test cases run")
	flag.Parse()
}

func main() {
	err := config.DefConfig.Init(TestConfig)
	if err != nil {
		log.Errorf("DefConfig.Init error:%s", err)
		return
	}

	rcSdk := poly_go_sdk.NewPolySdk()
	if err = btc.SetUpPoly(rcSdk, config.DefConfig.RchainJsonRpcAddress); err != nil {
		panic(err)
	}

	ethInvoker := eth.NewEInvoker()
	btcInvoker, err := btc.NewBtcInvoker(config.DefConfig.RchainJsonRpcAddress, config.DefConfig.RCWallet,
		config.DefConfig.RCWalletPwd, config.DefConfig.BtcRestAddr, config.DefConfig.BtcRestUser,
		config.DefConfig.BtcRestPwd, config.DefConfig.BtcSignerPrivateKey)
	if err != nil {
		panic(err)
	}
	ontInvoker, err := ont.NewOntInvoker(config.DefConfig.OntJsonRpcAddress, config.DefConfig.OntContractsAvmPath,
		config.DefConfig.OntWallet, config.DefConfig.OntWalletPassword)
	if err != nil {
		panic(err)
	}
	cmInvoker, err := cosmos.NewCosmosInvoker()
	if err != nil {
		panic(err)
	}

	testCases := make([]string, 0)
	if TestCases != "" {
		testCases = strings.Split(TestCases, ",")
	}
	testframework.TFramework.SetRcSdk(rcSdk)
	testframework.TFramework.SetEthInvoker(ethInvoker)
	testframework.TFramework.SetBtcInvoker(btcInvoker)
	testframework.TFramework.SetOntInvoker(ontInvoker)
	testframework.TFramework.SetCosmosInvoker(cmInvoker)

	//Start run test case
	testframework.TFramework.Run(testCases, LoopNumber)
	waitToExit()
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			fmt.Println("cross chain test received exit signal: ", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}
