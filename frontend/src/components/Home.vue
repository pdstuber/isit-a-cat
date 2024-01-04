<template>
    <div class="container">
        <div id="upload" v-if="isInitial || isSaving">
            <div class="file has-name is-centered">
                <label class="file-label" v-on:change="onFileChange">
                    <input class="file-input" type="file" ref="file" name="resume">
                    <span class="file-cta">
						<span class="file-icon">
							<i class="fas fa-upload"></i>
						</span>
						<span class="file-label">
							Choose a file…
						</span>
					</span>
                </label>
            </div>
        </div>
        <div v-if="isSuccess">
            <div v-if="prediction" class="box">
                <h2>I'm quite sure that this is <b>{{prediction.class !== 'cats' ? 'not' : ''}}</b> a cat</h2>
            </div>
            <div v-else class="box">
                <h2>Making prediction, please wait...</h2>
            </div>
            <div v-if="previewImageUrl">
                <img :src="previewImageUrl" class="img-responsive img-thumbnail" :alt="previewImageUrl.name">
            </div>
        </div>
        <div v-if="isFailed">
            <h2>Upload failed.</h2>
            <div class="notification">
                <pre>{{error}}</pre>
            </div>
        </div>
    </div>
</template>

<script>
    import axios from 'axios';
    import Pica from 'pica';

    const STATUS_INITIAL = 0, STATUS_SAVING = 1, STATUS_SUCCESS = 2, STATUS_FAILED = 3;
    const BASE_URL = 'localhost:8095';
    const TARGET_IMAGE_SIZE = 256
    const TARGET_IMAGE_MIME = 'image/jpeg'
    const TARGET_IMAGE_QUALITY = 0.99

    export default {
        name: 'Home',

        data() {
            return {
                error: null,
                currentStatus: STATUS_INITIAL,
                file: '',
                fileName: null,
                previewImageUrl: null,
                prediction: null,
                websocket: null
            }
        },
        computed: {
            isInitial() {
                return this.currentStatus === STATUS_INITIAL;
            },
            isSaving() {
                return this.currentStatus === STATUS_SAVING;
            },
            isSuccess() {
                return this.currentStatus === STATUS_SUCCESS;
            },
            isFailed() {
                return this.currentStatus === STATUS_FAILED;
            }
        },
        methods: {

            onFileChange(e) {
                let fileHandle = e.target.files[0];
                if (fileHandle.size > 0) {
                    this.handleFileUpload(fileHandle)
                }
            },

            async handleFileUpload(file) {
                this.currentStatus = STATUS_SAVING;
                const url = `http://${BASE_URL}/images`;
                this.previewImageUrl = await this.readFileAsync(file);
                let processedImage = await this.processImage(this.previewImageUrl);
                let formData = new FormData();
                formData.append('file', processedImage);
                axios.post(url,
                    formData,
                    {
                        headers: {
                            'Content-Type': 'multipart/form-data',
                        }
                    }
                ).then(x => {
                    this.currentStatus = STATUS_SUCCESS;
                    this.registerWebsocketForPredictionResult(x.data.id)
                })
                    .catch((error) => {
                        console.log('error in image upload');
                        this.currentStatus = STATUS_FAILED;
                        this.error = error
                    });
            },

            registerWebsocketForPredictionResult(id) {
                const url = `ws://${BASE_URL}/predictions/${id}`;
                this.websocket = new WebSocket(url);

                this.websocket.onerror = errorEvent => {
                    console.log(errorEvent)
                    this.currentStatus = STATUS_FAILED;
                    this.error = "server error"
                };

                this.websocket.onopen = openEvent => {
                    console.log("Websocket connection opened")
                    console.log(openEvent)
                }
                this.websocket.onclose = closeEvent => {
                    console.log("Websocket connection closed")
                    console.log(closeEvent)
                }

                this.websocket.onmessage = messageEvent => {
                    const wsMsg = messageEvent.data;
                    if (wsMsg.indexOf("error") > 0) {
                        console.error(wsMsg.error);
                    } else {
                        const predictionResult = JSON.parse(wsMsg)
                        if ('errorType' in predictionResult) {
                            console.log('error getting prediction result');
                            this.currentStatus = STATUS_FAILED;
                            this.error = predictionResult.message
                        } else {
                            this.prediction = predictionResult
                        }
                    }
                };
            },

            processImage(file) {
                return new Promise((resolve, reject) => {
                    const img = document.createElement('img');
                    img.onerror = reject;
                    img.onload = () => {
                        resolve(this.resizeAndCompress(img));
                    };

                    img.src = file;
                });
            },

            readFileAsync(file) {
                return new Promise((resolve, reject) => {
                    let reader = new FileReader();

                    reader.onload = (e) => {
                        resolve(e.target.result);
                    };

                    reader.onerror = reject;

                    if (file) {
                        if (/\.(jpe?g|png|gif)$/i.test(file.name)) {
                            reader.readAsDataURL(file);
                        }
                    }
                })
            },

            async resizeAndCompress(img) {
                const pica = Pica();
                const resizeCanvas = document.createElement('canvas');
                resizeCanvas.height = TARGET_IMAGE_SIZE;
                resizeCanvas.width = TARGET_IMAGE_SIZE;
                return await pica.resize(img, resizeCanvas)
                    .then(result => pica.toBlob(result, TARGET_IMAGE_MIME, TARGET_IMAGE_QUALITY))
                    .catch(error => {
                        console.log("error resizing image");
                        console.log(error);
                    });
            },
        }
    }
</script>

<style scoped>
    .file {
        margin-bottom: 5px;
    }
</style>
