/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package consumer

import (
	"time"

	"github.com/apache/rocketmq-client-go/internal"
	"github.com/apache/rocketmq-client-go/primitive"
)

type consumerOptions struct {
	internal.ClientOptions

	/**
	 * Backtracking consumption time with second precision. Time format is
	 * 20131223171201<br>
	 * Implying Seventeen twelve and 01 seconds on December 23, 2013 year<br>
	 * Default backtracking consumption time Half an hour ago.
	 */
	ConsumeTimestamp string

	// The socket timeout in milliseconds
	ConsumerPullTimeout time.Duration

	// Concurrently max span offset.it has no effect on sequential consumption
	ConsumeConcurrentlyMaxSpan int

	// Flow control threshold on queue level, each message queue will cache at most 1000 messages by default,
	// Consider the {PullBatchSize}, the instantaneous value may exceed the limit
	PullThresholdForQueue int64

	// Limit the cached message size on queue level, each message queue will cache at most 100 MiB messages by default,
	// Consider the {@code pullBatchSize}, the instantaneous value may exceed the limit
	//
	// The size of a message only measured by message body, so it's not accurate
	PullThresholdSizeForQueue int

	// Flow control threshold on topic level, default value is -1(Unlimited)
	//
	// The value of {@code pullThresholdForQueue} will be overwrote and calculated based on
	// {@code pullThresholdForTopic} if it is't unlimited
	//
	// For example, if the value of pullThresholdForTopic is 1000 and 10 message queues are assigned to this consumer,
	// then pullThresholdForQueue will be set to 100
	PullThresholdForTopic int

	// Limit the cached message size on topic level, default value is -1 MiB(Unlimited)
	//
	// The value of {@code pullThresholdSizeForQueue} will be overwrote and calculated based on
	// {@code pullThresholdSizeForTopic} if it is't unlimited
	//
	// For example, if the value of pullThresholdSizeForTopic is 1000 MiB and 10 message queues are
	// assigned to this consumer, then pullThresholdSizeForQueue will be set to 100 MiB
	PullThresholdSizeForTopic int

	// Message pull Interval
	PullInterval time.Duration

	// Batch consumption size
	ConsumeMessageBatchMaxSize int

	// Batch pull size
	PullBatchSize int32

	// Whether update subscription relationship when every pull
	PostSubscriptionWhenPull bool

	// Max re-consume times. -1 means 16 times.
	//
	// If messages are re-consumed more than {@link #maxReconsumeTimes} before Success, it's be directed to a deletion
	// queue waiting.
	MaxReconsumeTimes int

	// Suspending pulling time for cases requiring slow pulling like flow-control scenario.
	SuspendCurrentQueueTimeMillis time.Duration

	// Maximum amount of time a message may block the consuming thread.
	ConsumeTimeout time.Duration

	ConsumerModel  MessageModel
	Strategy       AllocateStrategy
	ConsumeOrderly bool
	FromWhere      ConsumeFromWhere

	Interceptors []primitive.Interceptor
	// TODO traceDispatcher
}

func defaultPushConsumerOptions() consumerOptions {
	opts := consumerOptions{
		ClientOptions: internal.DefaultClientOptions(),
		Strategy:      AllocateByAveragely,
	}
	opts.ClientOptions.GroupName = "DEFAULT_CONSUMER"
	return opts
}

type Option func(*consumerOptions)

func WithConsumerModel(m MessageModel) Option {
	return func(options *consumerOptions) {
		options.ConsumerModel = m
	}
}

func WithConsumeFromWhere(w ConsumeFromWhere) Option {
	return func(options *consumerOptions) {
		options.FromWhere = w
	}
}

// WithChainConsumerInterceptor returns a ConsumerOption that specifies the chained interceptor for consumer.
// The first interceptor will be the outer most, while the last interceptor will be the inner most wrapper
// around the real call.
func WithInterceptor(fs ...primitive.Interceptor) Option {
	return func(options *consumerOptions) {
		options.Interceptors = append(options.Interceptors, fs...)
	}
}

// WithGroupName set group name address
func WithGroupName(group string) Option {
	return func(opts *consumerOptions) {
		if group == "" {
			return
		}
		opts.GroupName = group
	}
}

// WithNameServer set NameServer address, only support one NameServer cluster in alpha2
func WithNameServer(nameServers []string) Option {
	return func(opts *consumerOptions) {
		if len(nameServers) > 0 {
			opts.NameServerAddrs = nameServers
		}
	}
}

// WithACL on/off ACL
func WithVIPChannel(enable bool) Option {
	return func(opts *consumerOptions) {
		opts.VIPChannelEnabled = enable
	}
}

// WithACL on/off ACL
func WithACL(enable bool) Option {
	return func(opts *consumerOptions) {
		opts.ACLEnabled = enable
	}
}

// WithRetry return a Option that specifies the retry times when send failed.
// TODO: use retry middleware instead
func WithRetry(retries int) Option {
	return func(opts *consumerOptions) {
		opts.RetryTimes = retries
	}
}
