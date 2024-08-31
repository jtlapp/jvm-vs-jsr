package com.joelapp.javabenchmarks.apiserver.controllers;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import com.joelapp.javabenchmarks.apiserver.models.Message;
import com.joelapp.javabenchmarks.apiserver.repositories.MessageRepo;

@RestController
@RequestMapping("/api")
public class MessageController {

    @Autowired
    private MessageRepo messageRepo;

    @PostMapping("/message")
    public Long createMessage(@RequestBody String text) {
        Message message = new Message();
        message.setText(text);
        message = messageRepo.save(message);
        return message.getId();
    }

    @GetMapping("/message/{id}")
    public String getMessage(@PathVariable Long id) {
        return messageRepo.findById(id)
                                .map(Message::getText)
                                .orElse("Not found");
    }
}
