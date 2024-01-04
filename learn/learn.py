#!/usr/bin/env python3

from PIL import ImageFile
from keras import models
from keras.applications.vgg16 import VGG16
from keras.layers import Dense, Dropout
from keras.layers import Flatten
from keras.optimizers import SGD
from keras.preprocessing import image

import export

ImageFile.LOAD_TRUNCATED_IMAGES = True
IMAGE_SIZE = 256


# define cnn model
def create_model():
    # load model
    model = VGG16(include_top=False, input_shape=(IMAGE_SIZE, IMAGE_SIZE, 3))
    # mark loaded layers as not trainable
    for layer in model.layers:
        layer.trainable = False
    # add new classifier layers
    x = Flatten()(model.layers[-1].output)

    x = Dense(256, activation='relu', kernel_initializer='he_uniform')(x)
    x = Dropout(0.2)(x)
    x = Dense(128, activation='relu', kernel_initializer='he_uniform')(x)
    output = Dense(2, activation='softmax')(x)

    # define new model
    model = models.Model(inputs=model.inputs, outputs=output)
    # compile model
    opt = SGD(lr=0.001, momentum=0.9)
    model.compile(loss='binary_crossentropy', optimizer=opt, metrics=['accuracy'])

    return model


def create_datagen():
    global datagen
    # create data generator
    datagen = image.ImageDataGenerator(
        featurewise_center=True,
        validation_split=0.2
    )
    # specify imagenet mean values for centering
    datagen.mean = [123.68, 116.779, 103.939]

    return datagen


def create_data_set(datagen, subset):
    return datagen.flow_from_directory(
        training_data_path,
        batch_size=128,
        target_size=(IMAGE_SIZE, IMAGE_SIZE),
        class_mode='categorical',
        subset=subset,
    )


training_data_path = 'training-images'
exported_model_folder = 'keras-exported-model'
model_export_path = f"{exported_model_folder}/model.h5"
labelsFilePath = f"{exported_model_folder}/labels.csv"
pb_model_name = 'model.pb'

classifier = create_model()

datagen = create_datagen()

training_data_set = create_data_set(datagen, 'training')
validation_data_set = create_data_set(datagen, 'validation')

history = classifier.fit_generator(
    training_data_set,
    steps_per_epoch=len(training_data_set),
    validation_data=validation_data_set,
    validation_steps=len(validation_data_set),
    epochs=10,
    verbose=1
)

_, acc = classifier.evaluate_generator(
    training_data_set, steps=len(training_data_set), verbose=0
)

print('Accuracy > %.3f' % (acc * 100.0))

classifier.save(model_export_path)

label_map = training_data_set.class_indices

with open(labelsFilePath, 'w') as f:
    for key in label_map.keys():
        f.write("%s,%s\n" % (label_map[key], key))

export.export_in_tf_format([out.op.name for out in classifier.outputs], exported_model_folder, pb_model_name)
