package com.jtlapp.jvmvsjs.joobymvc.controllers;

import io.jooby.annotation.*;

@Path("/")
public class HomeController {

  @GET
  public String home() {
    return "Running Jooby with MVC";
  }
}
