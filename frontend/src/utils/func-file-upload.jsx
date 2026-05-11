import {ArchiveOutline} from '@vicons/ionicons5'
import {ref} from "vue";
import _ from 'lodash-es'

/**
 * 文件上传组件
 * @param {Object} options - 配置选项
 * @param {Object} options.modal - naive-ui modal实例
 * @param {Object} options.loadingBar - naive-ui loadingBar实例
 * @param {Object} [options.multiple] - 默认false 是否支持多个文件 多选
 * @param {Object} [options.max] - naive-ui 默认undefined(不做限制) 限制上传文件数量
 * @param {Function} options.uploadFunc - 上传方法，点击上传后会传入文件信息file
 * @param {Object} [options.emits] - 事件发射器
 * @param {Function} [options.emits.uploadSuccess] - 事件发射器，上传成功防抖回调，会在所有文件上传后执行
 * @param {Object} [options.otherNUploadProps] - 其它n-upload属性配置 ，通过v-bind绑定
 * @returns {ModalReactive}
 */
export default (options) => {
    const {modal, loadingBar, emits, uploadFunc, multiple = false, max, otherNUploadProps = {}} = options
    const fileList = ref([])             // 文件列表
    const uploadRef = ref(null)          // 上传组件引用

    /* 上传逻辑 */
    async function customUpload(options) {
        try {
            await uploadFunc(options)
            handleUploadComplete()
        } catch (error) {
            loadingBar.error()
        }
    }

    function handleUploadComplete() {
        try {
            $message.success('上传成功')
            resetUploadState()
            notifyUploadSuccess()
        } finally {
            loadingBar.finish()
        }
    }

    function resetUploadState() {
        // 重置上传状态
        uploadRef.value?.clear()
        fileList.value = []
    }

    const notifyUploadSuccess = _.debounce(() => {
        emits?.uploadSuccess?.()
    }, 1000)

    /* 渲染模态框 */
    return modal.create({
        title: '上传文件',
        preset: 'card',
        style: {width: '600px', padding: '1rem'},
        maskClosable: false,
        content: () => (
            <div
                style="height: 100%; display: flex; flex-direction: column; gap: 1rem"
            >
                <n-upload multiple={multiple} max={max} {...otherNUploadProps}
                          file-list={fileList.value}
                          ref={uploadRef}
                          custom-request={customUpload}
                          default-upload={false}
                          onChange={(options) => {
                              // options 包含之前的文件和本次上传的文件
                              fileList.value = options.fileList
                          }}
                          directory-dnd
                >
                    <n-upload-dragger>
                        <div style="margin-bottom: 12px">
                            <n-icon size="48" depth={3}>
                                <ArchiveOutline/>
                            </n-icon>
                        </div>
                        <n-text style="font-size: 16px">
                            点击 | 拖动文件到该区域来上传
                        </n-text>
                    </n-upload-dragger>
                </n-upload>
                <n-button
                    disabled={fileList.value.length === 0}
                    onClick={() => uploadRef.value?.submit()}
                >
                    上传
                </n-button>
            </div>
        )
    })
}