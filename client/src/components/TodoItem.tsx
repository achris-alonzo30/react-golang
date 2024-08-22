import { Badge, Box, Flex, Spinner, Text, useToast } from "@chakra-ui/react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { FaCheckCircle } from "react-icons/fa";
import { MdDelete } from "react-icons/md";
import { Todo } from "./TodoList";
import { BASE_URL } from "../constants";

interface TodoItemProps {
	key: number;
	todo: Todo;
}

export const TodoItem = ({ todo }: TodoItemProps) => {
	const toast = useToast()
	const queryClient = useQueryClient()
	const { mutate: handleUpdateTodo, isPending: isUpdating } = useMutation({
		mutationKey: ["updateTodo"],
		mutationFn: async () => {
			try {
				const res = await fetch(`${BASE_URL}/api/todos/${todo._id}`, {
					method: "PATCH",
					headers: {
						"Content-Type": "application/json"
					},
				});

				if (!res.ok) {
					throw new Error("Failed to update todos");
				}

				const data = await res.json();

				return data;
			} catch (error) {
				throw new Error("Failed to update todos");
			}
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["todos"] })
		},
		onError: () => {
			toast({
				title: "Failed to update todo",
				status: "error",
				duration: 3000,
				isClosable: true,
			})
		}
	})

	const { mutate: handleDeleteTodo, isPending: isDeleting } = useMutation({
		mutationKey: ["deleteTodo"],
		mutationFn: async () => {
			try {
				const res = await fetch(`${BASE_URL}/api/todos/${todo._id}`, {
					method: "DELETE",
				});

				if (!res.ok) {
					throw new Error("Failed to delete todos");
				}

				const data = await res.json();

				return data;
			} catch (error) {
				throw new Error("Failed to delete todos");
			}
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["todos"] })
		},
		onError: () => {
			toast({
				title: "Failed to delete todo",
				status: "error",
				duration: 3000,
				isClosable: true,
			})
		}
	})


	return (
		<Flex gap={2} alignItems={"center"}>
			<Flex
				flex={1}
				alignItems={"center"}
				border={"1px"}
				borderColor={"gray.600"}
				p={2}
				borderRadius={"lg"}
				justifyContent={"space-between"}
			>
				<Text
					color={todo.completed ? "green.500" : "black.500"}
					textDecoration={todo.completed ? "line-through" : "none"}
				>
					{todo.title}
				</Text>
				{todo.completed && (
					<Badge ml='1' colorScheme='green'>
						Done
					</Badge>
				)}
				{!todo.completed && (
					<Badge ml='1' colorScheme='yellow'>
						In Progress
					</Badge>
				)}
			</Flex>
			<Flex gap={2} alignItems={"center"}>
				<Box color={"green.500"} cursor={"pointer"} onClick={() => handleUpdateTodo()}>
					{!isUpdating && <FaCheckCircle size={20} />}
					{isUpdating && <Spinner size="sm" />}
				</Box>
				<Box color={"red.500"} cursor={"pointer"} onClick={() => handleDeleteTodo()}>
					{!isDeleting && <MdDelete size={25} />}
					{isDeleting && <Spinner size="sm" />}
				</Box>
			</Flex>
		</Flex>
	);
};