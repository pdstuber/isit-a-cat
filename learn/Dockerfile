#FROM tensorflow/tensorflow:devel-gpu

FROM tensorflow/tensorflow

ADD train.py .

ARG TRAINING_DATA_LOCAL_PATH
RUN [ ! -z "${TRAINING_DATA_LOCAL_PATH}"] 
RUN mkdir keras-exported-model/

ENTRYPOINT [ "train.py" ]
