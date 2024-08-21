import {
  Box,
  Flex,
  Text,
  Button,
  Container,
  useColorMode,
  useColorModeValue,
} from "@chakra-ui/react";
import { LuSun } from "react-icons/lu";
import { IoMoon } from "react-icons/io5";

export const Navbar = () => {
  const { colorMode, toggleColorMode } = useColorMode();

  return (

    <Container maxW={"900px"}>
      <Box bg={useColorModeValue("black.500", "black.700")} px={4} my={4} borderRadius={"5"}>
        <Flex h={16} alignItems={"center"} justifyContent={"space-between"}>
          {/* LEFT SIDE */}
          <Flex
            justifyContent={"center"}
            alignItems={"center"}
            gap={3}
            display={{ base: "none", sm: "flex" }}
          >
            <Text fontSize={"40"}>React + Golang</Text>
          </Flex>

          {/* RIGHT SIDE */}
          <Flex alignItems={"center"} gap={3}>
            <Text fontSize={"lg"} fontWeight={500}>
              Daily Tasks
            </Text>
            {/* Toggle Color Mode */}
            <Button onClick={toggleColorMode}>
              {colorMode === "light" ? <IoMoon /> : <LuSun size={20} />}
            </Button>
          </Flex>
        </Flex>
      </Box>
    </Container>
  )
}