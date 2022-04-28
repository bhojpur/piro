FROM moby/buildkit:v0.9.3
WORKDIR /piro
COPY piro README.md /piro/
ENV PATH=/piro:$PATH
ENTRYPOINT [ "/bhojpur/piro" ]