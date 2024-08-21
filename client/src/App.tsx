import { Container, Stack } from "@chakra-ui/react";
import { Navbar } from "./components/Navbar";
import { TodoForm } from "./components/TodoForm";
import { TodoList } from "./components/TodoList";

export default function App() {
  return (
    <Stack h="100vh" w="100vw">
      <Navbar />
      <Container>
        <TodoForm />
        <TodoList />
      </Container>
    </Stack> 
  )
}