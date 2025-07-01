CC = gcc
CFLAGS = -Wall -O2 -std=c99
LIBS = -lGL -lGLU -lglut -lm
LIBS_PNG = -lm
TARGET = algebraic
TARGET_PNG = algebraic_png
TARGET_GO = algebraic_go
SOURCE = c.c
SOURCE_PNG = png_version.c
SOURCE_GO = algebraic.go

$(TARGET): $(SOURCE)
	$(CC) $(CFLAGS) -o $(TARGET) $(SOURCE) $(LIBS)

$(TARGET_PNG): $(SOURCE_PNG)
	$(CC) $(CFLAGS) -o $(TARGET_PNG) $(SOURCE_PNG) $(LIBS_PNG)

$(TARGET_GO): $(SOURCE_GO)
	go build -o $(TARGET_GO) $(SOURCE_GO)

png: $(TARGET_PNG)

go: $(TARGET_GO)

all: $(TARGET) $(TARGET_PNG) $(TARGET_GO)

clean:
	rm -f $(TARGET) $(TARGET_PNG) $(TARGET_GO) *.log *.csv *.txt *.ppm *.png

.PHONY: clean png go all