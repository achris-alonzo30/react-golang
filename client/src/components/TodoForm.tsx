import { 
  Flex,
  Input,
  Button,  
  Spinner, 
  useToast
} from "@chakra-ui/react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { IoMdAdd } from "react-icons/io";
import { BASE_URL } from "../constants";

export const TodoForm = () => {
	const toast = useToast();
	const queryClient = useQueryClient();
	const [newTodo, setNewTodo] = useState("");

	const { mutate: createTodo, isPending } = useMutation({
		mutationKey: ["createTodo"],
		mutationFn: async (e: React.FormEvent) => {
			e.preventDefault();
			try {
				const res = await fetch(`${BASE_URL}/api/todos`, {
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					body: JSON.stringify({ title: newTodo }),
				});

				if (!res.ok) {
					throw new Error("Failed to create todo");
				}
	
				const data = await res.json();
				setNewTodo("");
				return data;
			} catch (error) {
				throw new Error("Failed to create todo");
			}
		},
		onSuccess: () => {
			toast({
				title: "Todo created",
				status: "success",
				duration: 3000,
				isClosable: true,
			})
			queryClient.invalidateQueries({ queryKey: ["todos"] });
		},
		onError: () => {
			toast({
				title: "Failed to create todo",
				status: "error",
				duration: 3000,
				isClosable: true,
			})
		}
	})
  
	return (
		<form onSubmit={createTodo}>
			<Flex gap={2}>
				<Input
					type='text'
					value={newTodo}
					onChange={(e) => setNewTodo(e.target.value)}
					ref={(input) => input && input.focus()}
				/>
				<Button
					mx={2}
					type='submit'
					_active={{
						transform: "scale(.97)",
					}}
				>
					{isPending ? <Spinner size={"xs"} /> : <IoMdAdd size={30} />}
				</Button>
			</Flex>
		</form>
	);
};