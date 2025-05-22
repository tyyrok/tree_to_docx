from pydantic import EmailStr
from pydantic.networks import validate_email
from pydantic_core import core_schema


class EmailStrLower(EmailStr):
    @classmethod
    def _validate(
        cls, __input_value: str, _: core_schema.ValidationInfo  # noqa: PYI063
    ) -> str:
        return validate_email(__input_value)[1].lower()
