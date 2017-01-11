FROM iron/go
ADD server /

ENTRYPOINT ["./server"]



