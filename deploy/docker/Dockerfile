FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
COPY app.jar app.jar
EXPOSE 8443
ENTRYPOINT ["java", "-jar", "app.jar"]
