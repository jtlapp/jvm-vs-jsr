FROM openjdk:21-jdk-slim
VOLUME /tmp
COPY target/spring-kernel-threads-0.1.0-SNAPSHOT.jar app.jar
ENTRYPOINT ["java", "-jar", "/app.jar"]
