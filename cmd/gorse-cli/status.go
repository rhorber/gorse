// Copyright 2021 gorse Project Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/zhenghaoz/gorse/base"
	"github.com/zhenghaoz/gorse/protocol"
	"github.com/zhenghaoz/gorse/storage/cache"
	"go.uber.org/zap"
	"os"
)

func init() {
	cliCommand.AddCommand(clusterCommand)
	cliCommand.AddCommand(statusCommand)
	cliCommand.AddCommand(configCommand)
}

var clusterCommand = &cobra.Command{
	Use:   "cluster",
	Short: "cluster information",
	Run: func(cmd *cobra.Command, args []string) {
		cluster, err := masterClient.GetMeta(context.Background(), &protocol.RequestInfo{})
		if err != nil {
			base.Logger().Fatal("failed to get cluster information", zap.Error(err))
		}
		// show cluster
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"role", "address"})
		for _, addr := range cluster.Servers {
			table.Append([]string{"server", addr})
		}
		for _, addr := range cluster.Workers {
			table.Append([]string{"worker", addr})
		}
		table.Render()
	},
}

var statusCommand = &cobra.Command{
	Use:   "status",
	Short: "status of recommender system",
	Run: func(cmd *cobra.Command, args []string) {
		// connect to cache store
		cacheStore, err := cache.Open(globalConfig.Database.CacheStore)
		if err != nil {
			base.Logger().Fatal("failed to connect database", zap.Error(err))
		}
		// show status
		status := []string{
			cache.CollectPopularTime,
			cache.CollectLatestTime,
			cache.CollectSimilarTime,
			cache.FitMatrixFactorizationTime,
			cache.FitFactorizationMachineTime,
			cache.MatrixFactorizationVersion,
			cache.FactorizationMachineVersion,
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"status", "value"})
		for _, stat := range status {
			val, err := cacheStore.GetString(cache.GlobalMeta, stat)
			if err != nil && err.Error() != "redis: nil" {
				base.Logger().Fatal("failed to get meta", zap.Error(err))
			}
			table.Append([]string{stat, val})
		}
		table.Render()
	},
}

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "config of recommender system",
	Run: func(cmd *cobra.Command, args []string) {
		bytes, err := json.MarshalIndent(globalConfig, "", "\t")
		if err != nil {
			base.Logger().Fatal("failed to marshall JSON", zap.Error(err))
		}
		fmt.Println(string(bytes))
	},
}
