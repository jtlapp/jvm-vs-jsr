package com.joelapp.javabenchmarks.hikariquerylibrary.tables;

abstract public class Table {
    private static final int ZERO_PADDING_WIDTH = 10;

    static String createID(String prefix, int value) {
        return prefix + String.format("%1$" + ZERO_PADDING_WIDTH + "s", value)
                .replace(' ', '0');
    }
}
