package com.jtlapp.jvmvsjs.controllers;

import io.jooby.annotation.*;

@Path("/")
public class Controller {

  @GET
  public String sayHi() {
    return "Welcome to Jooby!";
  }
}
