/*
 * Copyright © 2021 - 2022 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package job

// InitJobService 启动JobTask服务
func InitJobService() error {
	go Executor.Start()
	go Fiexd.Start()
	select {}
}
