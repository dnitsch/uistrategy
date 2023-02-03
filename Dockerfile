FROM debian:bullseye

WORKDIR /seeder 

COPY dist/uiseeder-linux uiseeder
COPY .ignore-paypal.yml paypal.yml

# RUN apt-get update \
# 	&& apt install -y libnss3

ENTRYPOINT [ "/seeder/uiseeder", "-i", "/seeder/paypal.yml" ]