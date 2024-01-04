#!/usr/bin/env python3
import keras
from PIL import Image
from keras.engine.saving import load_model

# load and prepare the image
from keras_applications.vgg16 import preprocess_input
from keras_preprocessing import image
from np.magic import np

import export

IMG_SIZE = 256


def load_image(filename):
    # load the image
    img = image.load_img(filename, target_size=(IMG_SIZE, IMG_SIZE))
    # convert to array
    img = image.img_to_array(img)
    # reshape into a single sample with 3 channels
    img = img.reshape(1, IMG_SIZE, IMG_SIZE, 3)
    # center pixel data
    img = img.astype('float32')
    img = img - [123.68, 116.779, 103.939]
    return img


model_export_path = 'keras-exported-model/model.h5'
keras.backend.set_learning_phase(0)
classifier = load_model(model_export_path)
print("### Inputs ###")
print(classifier.inputs)
print("### Outputs ###")
print(classifier.outputs)
print(classifier.summary())

img = load_image('test-images/cats/cat.11446.jpg')
# predict the class
result = classifier.predict(img)

# img = load_image('test-images/non_cats/lion.jpeg')
# predict the class
# result = classifier.predict(img, steps=1)
print(f"[{result[0][0]:.2f}][{result[0][1]:.2f}]")
print([out.op.name for out in classifier.outputs])
export.export_in_tf_format([out.op.name for out in classifier.outputs], 'keras-exported-model', 'out.pb')

