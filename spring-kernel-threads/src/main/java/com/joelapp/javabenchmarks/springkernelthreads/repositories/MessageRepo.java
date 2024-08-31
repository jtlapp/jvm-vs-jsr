package com.joelapp.javabenchmarks.apiserver.repositories;

import org.springframework.data.jpa.repository.JpaRepository;

import com.joelapp.javabenchmarks.apiserver.models.Message;

public interface MessageRepo extends JpaRepository<Message, Long> {
}
