/*
 * Copyright © 2021 - 2023 vity <vityme@icloud.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 */

package job

// InitJobService 启动JobTask服务
func InitJobService() error {
	defer func() {
		select {}
	}()
	go Task.Start()
	go Fiexd.Start()
	return nil
}
